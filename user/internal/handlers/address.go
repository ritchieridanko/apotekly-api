package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/dto"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"go.opentelemetry.io/otel"
)

const addressErrorTracer string = "handler.address"

type AddressHandler interface {
	NewAddress(ctx *gin.Context)
	GetAllAddresses(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	DeleteAddress(ctx *gin.Context)
}

type addressHandler struct {
	au usecases.AddressUsecase
}

func NewAddressHandler(au usecases.AddressUsecase) AddressHandler {
	return &addressHandler{au}
}

func (h *addressHandler) NewAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "NewAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqNewAddress
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if validateErr := utils.ValidateNewAddress(payload); validateErr != "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, validateErr, errors.New("invalid payload"))
		ctx.Error(err)
		return
	}

	data := entities.NewAddress{
		Receiver:    strings.TrimSpace(payload.Receiver),
		Phone:       payload.Phone,
		Label:       strings.TrimSpace(payload.Label),
		Notes:       payload.Notes,
		IsPrimary:   payload.IsPrimary,
		Country:     strings.ToLower(strings.TrimSpace(payload.Country)),
		AdminLevel1: utils.NormalizePtr(payload.AdminLevel1),
		AdminLevel2: utils.NormalizePtr(payload.AdminLevel2),
		AdminLevel3: utils.NormalizePtr(payload.AdminLevel3),
		AdminLevel4: utils.NormalizePtr(payload.AdminLevel4),
		Street:      strings.TrimSpace(payload.Street),
		PostalCode:  strings.TrimSpace(payload.PostalCode),
		Latitude:    payload.Latitude,
		Longitude:   payload.Longitude,
	}

	address, unsetPrimaryID, err := h.au.NewAddress(ctxWithTracer, authID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	var unsetID *int64
	if unsetPrimaryID != 0 {
		unsetID = &unsetPrimaryID
	}

	response := dto.RespNewAddress{
		Created:        h.setAddressAsResponse(*address),
		UnsetPrimaryID: unsetID,
	}

	utils.SetResponse(ctx, "Address added successfully.", response, http.StatusCreated)
}

func (h *addressHandler) GetAllAddresses(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "GetAllAddresses")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	addresses, err := h.au.GetAllAddresses(ctxWithTracer, authID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := make([]dto.RespAddress, 0, len(addresses))
	for _, address := range addresses {
		addr := h.setAddressAsResponse(address)
		response = append(response, addr)
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *addressHandler) UpdateAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "UpdateAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	addressID, err := utils.ToInt64(ctx.Param("id"))
	if err != nil {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqUpdateAddress
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if validateErr := utils.ValidateAddressChange(payload); validateErr != "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, validateErr, errors.New("invalid payload"))
		ctx.Error(err)
		return
	}

	data := entities.AddressChange{
		Receiver:    utils.TrimSpacePtr(payload.Receiver),
		Phone:       payload.Phone,
		Label:       utils.TrimSpacePtr(payload.Label),
		Notes:       payload.Notes,
		IsPrimary:   payload.IsPrimary,
		Country:     utils.NormalizePtr(payload.Country),
		AdminLevel1: utils.NormalizePtr(payload.AdminLevel1),
		AdminLevel2: utils.NormalizePtr(payload.AdminLevel2),
		AdminLevel3: utils.NormalizePtr(payload.AdminLevel3),
		AdminLevel4: utils.NormalizePtr(payload.AdminLevel4),
		Street:      utils.TrimSpacePtr(payload.Street),
		PostalCode:  utils.TrimSpacePtr(payload.PostalCode),
		Latitude:    payload.Latitude,
		Longitude:   payload.Longitude,
	}

	address, unsetPrimaryID, err := h.au.UpdateAddress(ctxWithTracer, authID, addressID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	var unsetID *int64
	if unsetPrimaryID != 0 {
		unsetID = &unsetPrimaryID
	}

	response := dto.RespUpdateAddress{
		Updated:        h.setAddressAsResponse(*address),
		UnsetPrimaryID: unsetID,
	}

	utils.SetResponse(ctx, "Address updated successfully.", response, http.StatusOK)
}

func (h *addressHandler) DeleteAddress(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(addressErrorTracer).Start(ctx.Request.Context(), "DeleteAddress")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	addressID, err := utils.ToInt64(ctx.Param("id"))
	if err != nil {
		err := ce.NewError(span, ce.CodeInvalidParams, ce.MsgInvalidParams, err)
		ctx.Error(err)
		return
	}

	deletedID, newPrimaryID, err := h.au.DeleteAddress(ctxWithTracer, authID, addressID)
	if err != nil {
		ctx.Error(err)
		return
	}

	var updatedID *int64
	if newPrimaryID != 0 {
		updatedID = &newPrimaryID
	}

	response := dto.RespDeleteAddress{
		DeletedID:    deletedID,
		NewPrimaryID: updatedID,
	}

	utils.SetResponse(ctx, "Address deleted successfully.", response, http.StatusOK)
}

func (h *addressHandler) setAddressAsResponse(address entities.Address) dto.RespAddress {
	return dto.RespAddress{
		ID:          address.ID,
		Receiver:    address.Receiver,
		Phone:       address.Phone,
		Label:       address.Label,
		Notes:       address.Notes,
		IsPrimary:   address.IsPrimary,
		Country:     utils.ToTitlecase(address.Country),
		AdminLevel1: utils.ToTitlecasePtr(address.AdminLevel1),
		AdminLevel2: utils.ToTitlecasePtr(address.AdminLevel2),
		AdminLevel3: utils.ToTitlecasePtr(address.AdminLevel3),
		AdminLevel4: utils.ToTitlecasePtr(address.AdminLevel4),
		Street:      address.Street,
		PostalCode:  address.PostalCode,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
	}
}
