package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/dto"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"github.com/ritchieridanko/apotekly-api/user/pkg/ce"
)

// TODO
// (1): Implement Image Upload to Cloud

const UserErrorTracer = ce.UserHandlerTracer

type UserHandler interface {
	NewUser(ctx *gin.Context)
}

type userHandler struct {
	uu usecases.UserUsecase
}

func NewUserHandler(uu usecases.UserUsecase) UserHandler {
	return &userHandler{uu}
}

func (h *userHandler) NewUser(ctx *gin.Context) {
	tracer := UserErrorTracer + ": NewUser()"

	authID, err := utils.GetAuthIDFromContext(ctx)
	if err != nil {
		err := ce.NewError(ce.ErrCodeContext, ce.ErrMsgInternalServer, tracer, err)
		ctx.Error(err)
		return
	}

	var payload dto.ReqNewUser
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
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

	if err := h.uu.NewUser(ctx, authID, &data); err != nil {
		ctx.Error(err)
		return
	}

	utils.SetResponse(ctx, "user created successfully", nil, http.StatusCreated)
}
