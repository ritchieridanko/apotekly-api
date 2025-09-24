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

	if payload.AdminLevel1 != nil {
		normalized := strings.ToLower(strings.TrimSpace(*payload.AdminLevel1))
		payload.AdminLevel1 = &normalized
	}
	if payload.AdminLevel2 != nil {
		normalized := strings.ToLower(strings.TrimSpace(*payload.AdminLevel2))
		payload.AdminLevel2 = &normalized
	}
	if payload.AdminLevel3 != nil {
		normalized := strings.ToLower(strings.TrimSpace(*payload.AdminLevel3))
		payload.AdminLevel3 = &normalized
	}
	if payload.AdminLevel4 != nil {
		normalized := strings.ToLower(strings.TrimSpace(*payload.AdminLevel4))
		payload.AdminLevel4 = &normalized
	}

	data := entities.NewAddress{
		Receiver:    strings.TrimSpace(payload.Receiver),
		Phone:       payload.Phone,
		Label:       strings.TrimSpace(payload.Label),
		Notes:       payload.Notes,
		IsPrimary:   payload.IsPrimary,
		Country:     strings.ToLower(strings.TrimSpace(payload.Country)),
		AdminLevel1: payload.AdminLevel1,
		AdminLevel2: payload.AdminLevel2,
		AdminLevel3: payload.AdminLevel3,
		AdminLevel4: payload.AdminLevel4,
		Street:      strings.TrimSpace(payload.Street),
		PostalCode:  strings.TrimSpace(payload.PostalCode),
		Latitude:    payload.Latitude,
		Longitude:   payload.Longitude,
	}

	address, err := h.au.NewAddress(ctxWithTracer, authID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	if address.AdminLevel1 != nil {
		normalized := utils.ToTitlecase(*address.AdminLevel1)
		address.AdminLevel1 = &normalized
	}
	if address.AdminLevel2 != nil {
		normalized := utils.ToTitlecase(*address.AdminLevel2)
		address.AdminLevel2 = &normalized
	}
	if address.AdminLevel3 != nil {
		normalized := utils.ToTitlecase(*address.AdminLevel3)
		address.AdminLevel3 = &normalized
	}
	if address.AdminLevel4 != nil {
		normalized := utils.ToTitlecase(*address.AdminLevel4)
		address.AdminLevel4 = &normalized
	}

	response := dto.RespNewAddress{
		ID:          address.ID,
		Receiver:    address.Receiver,
		Phone:       address.Phone,
		Label:       address.Label,
		Notes:       address.Notes,
		IsPrimary:   address.IsPrimary,
		Country:     utils.ToTitlecase(address.Country),
		AdminLevel1: address.AdminLevel1,
		AdminLevel2: address.AdminLevel2,
		AdminLevel3: address.AdminLevel3,
		AdminLevel4: address.AdminLevel4,
		Street:      address.Street,
		PostalCode:  address.PostalCode,
		Latitude:    address.Latitude,
		Longitude:   address.Longitude,
	}

	utils.SetResponse(ctx, "Address added successfully.", response, http.StatusCreated)
}
