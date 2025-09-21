package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/dto"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"go.opentelemetry.io/otel"
)

// TODO
// (1): Implement image upload to cloud

const userErrorTracer string = "handler.user"

type UserHandler interface {
	NewUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
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

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		err := ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqNewUser
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
		ctx.Error(err)
		return
	}

	if !utils.ValidateNewUser(payload) {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, errors.New("invalid payload"))
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

	// TODO (1)

	user, err := h.uu.NewUser(ctxWithTracer, authID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespNewUser{
		UserID:         user.UserID,
		Name:           user.Name,
		Bio:            user.Bio,
		Sex:            user.Sex,
		Birthdate:      user.Birthdate,
		Phone:          user.Phone,
		ProfilePicture: user.ProfilePicture,
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

	response := dto.RespNewUser{
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
