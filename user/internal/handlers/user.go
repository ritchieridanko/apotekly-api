package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ritchieridanko/apotekly-api/user/config"
	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/constants"
	"github.com/ritchieridanko/apotekly-api/user/internal/dto"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "handler.user"

type UserHandler interface {
	NewUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	ChangeProfilePicture(ctx *gin.Context)
}

type userHandler struct {
	uu usecases.UserUsecase
}

func NewUserHandler(uu usecases.UserUsecase) UserHandler {
	return &userHandler{uu}
}

func (h *userHandler) NewUser(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "NewUser")
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

	var payload dto.ReqNewUser
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if validateErr := utils.ValidateNewUser(payload); validateErr != "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, validateErr, errors.New("invalid payload"))
		ctx.Error(err)
		return
	}

	data := entities.NewUser{
		Name:      payload.Name,
		Bio:       payload.Bio,
		Sex:       payload.Sex,
		Birthdate: payload.Birthdate,
		Phone:     payload.Phone,
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

	user, err := h.uu.NewUser(ctxWithTracer, authID, &data, image)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespNewUser{
		Created: dto.RespUser{
			UserID:         user.UserID,
			Name:           user.Name,
			Bio:            user.Bio,
			Sex:            user.Sex,
			Birthdate:      user.Birthdate,
			Phone:          user.Phone,
			ProfilePicture: user.ProfilePicture,
		},
	}

	utils.SetResponse(ctx, "User created successfully.", response, http.StatusCreated)
}

func (h *userHandler) GetUser(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "GetUser")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	user, err := h.uu.GetUser(ctxWithTracer, authID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespUser{
		UserID:         user.UserID,
		Name:           user.Name,
		Bio:            user.Bio,
		Sex:            user.Sex,
		Birthdate:      user.Birthdate,
		Phone:          user.Phone,
		ProfilePicture: user.ProfilePicture,
	}

	utils.SetResponse(ctx, "ok", response, http.StatusOK)
}

func (h *userHandler) UpdateUser(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "UpdateUser")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqUpdateUser
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if validateErr := utils.ValidateUserChange(payload); validateErr != "" {
		err := ce.NewError(span, ce.CodeInvalidPayload, validateErr, errors.New("invalid payload"))
		ctx.Error(err)
		return
	}

	data := entities.UserChange{
		Name:      payload.Name,
		Bio:       payload.Bio,
		Sex:       payload.Sex,
		Birthdate: payload.Birthdate,
		Phone:     payload.Phone,
	}

	user, err := h.uu.UpdateUser(ctxWithTracer, authID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespUpdateUser{
		Updated: dto.RespUser{
			UserID:         user.UserID,
			Name:           user.Name,
			Bio:            user.Bio,
			Sex:            user.Sex,
			Birthdate:      user.Birthdate,
			Phone:          user.Phone,
			ProfilePicture: user.ProfilePicture,
		},
	}

	utils.SetResponse(ctx, "User updated successfully.", response, http.StatusOK)
}

func (h *userHandler) ChangeProfilePicture(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "ChangeProfilePicture")
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

	file, err := ctx.FormFile("image")
	if err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	image, err := file.Open()
	if err != nil {
		err := ce.NewError(span, ce.CodeRequestFile, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}
	defer image.Close()

	if file.Size > maxSize {
		err := ce.NewError(span, ce.CodeInvalidPayload, "File size exceeds limit.", errors.New("file size exceeds maximum size"))
		ctx.Error(err)
		return
	}

	if err := h.uu.ChangeProfilePicture(ctxWithTracer, authID, image); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "Profile picture changed.", nil, http.StatusOK)
}
