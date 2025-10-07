package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/dto"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/validator"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/utils"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"go.opentelemetry.io/otel"
)

const addressErrorTracer string = "handler.address"

type AddressHandler struct {
	au        usecases.AddressUsecase
	validator *validator.Validator
}

func NewAddressHandler(au usecases.AddressUsecase, validator *validator.Validator) *AddressHandler {
	return &AddressHandler{au, validator}
}

func (h *AddressHandler) CreateAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "CreateAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to create address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	var payload dto.CreateAddressRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to create address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		wErr := fmt.Errorf("failed to create address: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, wErr))
		return
	}

	data := entities.CreateAddress{
		Recipient:    strings.TrimSpace(payload.Recipient),
		Phone:        strings.TrimSpace(payload.Phone),
		Label:        strings.TrimSpace(payload.Label),
		Notes:        payload.Notes,
		IsPrimary:    payload.IsPrimary,
		Country:      utils.Normalize(payload.Country),
		Subdivision1: utils.NormalizePtr(payload.Subdivision1),
		Subdivision2: utils.NormalizePtr(payload.Subdivision2),
		Subdivision3: utils.NormalizePtr(payload.Subdivision3),
		Subdivision4: utils.NormalizePtr(payload.Subdivision4),
		Street:       strings.TrimSpace(payload.Street),
		PostalCode:   strings.TrimSpace(payload.PostalCode),
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
	}

	address, oldPrimaryAddress, err := h.au.CreateAddress(ctxWithTracer, authID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	var updated *dto.AddressResponse
	if oldPrimaryAddress != nil {
		resp := h.addressToResponse(*oldPrimaryAddress)
		updated = &resp
	}

	response := dto.CreateAddressResponse{
		Created: h.addressToResponse(*address),
		Updated: updated,
	}

	utils.SetResponse(ctx, "Address created successfully", response, http.StatusCreated)
}

func (h *AddressHandler) GetAllAddresses(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "GetAllAddresses")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch all addresses: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	addresses, err := h.au.GetAllAddresses(ctxWithTracer, authID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := make([]dto.AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		addr := h.addressToResponse(address)
		response = append(response, addr)
	}

	utils.SetResponse(ctx, "Addresses retrieved successfully", response, http.StatusOK)
}

func (h *AddressHandler) UpdateAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "UpdateAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to update address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	addressID, err := utils.ToInt64(ctx.Param("id"))
	if err != nil {
		wErr := fmt.Errorf("failed to update address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	var payload dto.UpdateAddressRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to update address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		wErr := fmt.Errorf("failed to update address: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, wErr))
		return
	}

	data := entities.UpdateAddress{
		Recipient:    utils.TrimSpacePtr(payload.Recipient),
		Phone:        utils.TrimSpacePtr(payload.Phone),
		Label:        utils.TrimSpacePtr(payload.Label),
		Notes:        payload.Notes,
		Country:      utils.NormalizePtr(payload.Country),
		Subdivision1: utils.NormalizePtr(payload.Subdivision1),
		Subdivision2: utils.NormalizePtr(payload.Subdivision2),
		Subdivision3: utils.NormalizePtr(payload.Subdivision3),
		Subdivision4: utils.NormalizePtr(payload.Subdivision4),
		Street:       utils.TrimSpacePtr(payload.Street),
		PostalCode:   utils.TrimSpacePtr(payload.PostalCode),
		Latitude:     payload.Latitude,
		Longitude:    payload.Longitude,
	}

	address, err := h.au.UpdateAddress(ctxWithTracer, authID, addressID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.UpdateAddressResponse{
		Updated: h.addressToResponse(*address),
	}

	utils.SetResponse(ctx, "Address updated successfully", response, http.StatusOK)
}

func (h *AddressHandler) SetPrimaryAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "SetPrimaryAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to set primary address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	addressID, err := utils.ToInt64(ctx.Param("id"))
	if err != nil {
		wErr := fmt.Errorf("failed to set primary address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	newPrimaryAddress, oldPrimaryAddress, err := h.au.SetPrimaryAddress(ctxWithTracer, authID, addressID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.SetPrimaryAddressResponse{
		NewPrimaryAddress: h.addressToResponse(*newPrimaryAddress),
		OldPrimaryAddress: h.addressToResponse(*oldPrimaryAddress),
	}

	utils.SetResponse(ctx, "Address set primary successfully", response, http.StatusOK)
}

func (h *AddressHandler) DeleteAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "DeleteAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to delete address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	addressID, err := utils.ToInt64(ctx.Param("id"))
	if err != nil {
		wErr := fmt.Errorf("failed to delete address: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, wErr))
		return
	}

	newPrimaryAddress, err := h.au.DeleteAddress(ctxWithTracer, authID, addressID)
	if err != nil {
		ctx.Error(err)
		return
	}

	msg := "Address deleted successfully"
	if newPrimaryAddress != nil {
		response := dto.DeleteAddressResponse{
			NewPrimaryAddress: h.addressToResponse(*newPrimaryAddress),
		}

		utils.SetResponse(ctx, msg, response, http.StatusOK)
		return
	}

	utils.SetResponse(ctx, msg, nil, http.StatusOK)
}

func (h *AddressHandler) addressToResponse(address entities.Address) dto.AddressResponse {
	return dto.AddressResponse{
		ID:           address.ID,
		Recipient:    address.Recipient,
		Phone:        address.Phone,
		Label:        address.Label,
		Notes:        address.Notes,
		IsPrimary:    address.IsPrimary,
		Country:      utils.ToTitlecase(address.Country),
		Subdivision1: utils.ToTitlecasePtr(address.Subdivision1),
		Subdivision2: utils.ToTitlecasePtr(address.Subdivision2),
		Subdivision3: utils.ToTitlecasePtr(address.Subdivision3),
		Subdivision4: utils.ToTitlecasePtr(address.Subdivision4),
		Street:       address.Street,
		PostalCode:   address.PostalCode,
		Latitude:     address.Latitude,
		Longitude:    address.Longitude,
		CreatedAt:    address.CreatedAt,
		UpdatedAt:    address.UpdatedAt,
	}
}
