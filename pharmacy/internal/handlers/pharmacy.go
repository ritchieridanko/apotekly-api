package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ritchieridanko/apotekly-api/pharmacy/config"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/ce"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/constants"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/dto"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/entities"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/utils"
	"go.opentelemetry.io/otel"
)

const pharmacyErrorTracer string = "handler.pharmacy"

type PharmacyHandler interface {
	NewPharmacy(ctx *gin.Context)
	GetPharmacy(ctx *gin.Context)
}

type pharmacyHandler struct {
	pu usecases.PharmacyUsecase
}

func NewPharmacyHandler(pu usecases.PharmacyUsecase) PharmacyHandler {
	return &pharmacyHandler{pu}
}

func (h *pharmacyHandler) NewPharmacy(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(pharmacyErrorTracer).Start(ctx.Request.Context(), "NewPharmacy")
	defer span.End()

	// limit request body size to max size
	maxSize := constants.SizeMB * config.StorageGetImageMaxSize()
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqNewPharmacy
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if validateErr := utils.ValidateNewPharmacy(payload); validateErr != "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, validateErr, errors.New("invalid payload"))
		ctx.Error(err)
		return
	}

	data := entities.NewPharmacy{
		Name:             strings.TrimSpace(payload.Name),
		LegalName:        utils.TrimSpacePtr(payload.LegalName),
		Description:      payload.Description,
		LicenseNumber:    payload.LicenseNumber,
		LicenseAuthority: payload.LicenseAuthority,
		LicenseExpiry:    payload.LicenseExpiry,
		Email:            payload.Email,
		Phone:            payload.Phone,
		Website:          payload.Website,
		Country:          utils.Normalize(payload.Country),
		AdminLevel1:      utils.NormalizePtr(payload.AdminLevel1),
		AdminLevel2:      utils.NormalizePtr(payload.AdminLevel2),
		AdminLevel3:      utils.NormalizePtr(payload.AdminLevel3),
		AdminLevel4:      utils.NormalizePtr(payload.AdminLevel4),
		Street:           strings.TrimSpace(payload.Street),
		PostalCode:       strings.TrimSpace(payload.PostalCode),
		Latitude:         payload.Latitude,
		Longitude:        payload.Longitude,
		OpeningHours:     payload.OpeningHours,
	}

	var image multipart.File
	if payload.Image != nil {
		file, err := payload.Image.Open()
		if err != nil {
			log.Println("WARNING -> failed to open and upload image:", err.Error())
		} else {
			defer file.Close()

			if payload.Image.Size > maxSize {
				err := ce.NewError(span, ce.CodeInvalidPayload, "File size exceeds limit.", errors.New("file size exceeds maximum size"))
				ctx.Error(err)
				return
			}
			image = file
		}
	}

	pharmacy, err := h.pu.NewPharmacy(ctxWithTracer, authID, &data, image)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespNewPharmacy{
		Created: h.setPharmacyAsResponse(*pharmacy),
	}

	utils.SetResponse(ctx, "Pharmacy created successfully.", response, http.StatusCreated)
}

func (h *pharmacyHandler) GetPharmacy(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(pharmacyErrorTracer).Start(ctx.Request.Context(), "GetPharmacy")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	pharmacy, err := h.pu.GetPharmacy(ctxWithTracer, authID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := h.setPharmacyAsResponse(*pharmacy)

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *pharmacyHandler) setPharmacyAsResponse(pharmacy entities.Pharmacy) dto.RespPharmacy {
	return dto.RespPharmacy{
		ID:               pharmacy.PharmacyPublicID,
		Name:             pharmacy.Name,
		LegalName:        pharmacy.LegalName,
		Description:      pharmacy.Description,
		LicenseNumber:    pharmacy.LicenseNumber,
		LicenseAuthority: pharmacy.LicenseAuthority,
		LicenseExpiry:    pharmacy.LicenseExpiry,
		Email:            pharmacy.Email,
		Phone:            pharmacy.Phone,
		Website:          pharmacy.Website,
		Country:          utils.ToTitlecase(pharmacy.Country),
		AdminLevel1:      utils.ToTitlecasePtr(pharmacy.AdminLevel1),
		AdminLevel2:      utils.ToTitlecasePtr(pharmacy.AdminLevel2),
		AdminLevel3:      utils.ToTitlecasePtr(pharmacy.AdminLevel3),
		AdminLevel4:      utils.ToTitlecasePtr(pharmacy.AdminLevel4),
		Street:           pharmacy.Street,
		PostalCode:       pharmacy.PostalCode,
		Latitude:         pharmacy.Latitude,
		Longitude:        pharmacy.Longitude,
		Logo:             pharmacy.Logo,
		OpeningHours:     pharmacy.OpeningHours,
		Status:           pharmacy.Status,
	}
}
