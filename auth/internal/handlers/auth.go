package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dtos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const AuthErrorTracer = ce.AuthHandlerTracer

type AuthHandler interface {
	Register(ctx *gin.Context)
}

type authHandler struct {
	au usecases.AuthUsecase
}

func NewAuthHandler(au usecases.AuthUsecase) AuthHandler {
	return &authHandler{au}
}

func (h *authHandler) Register(ctx *gin.Context) {
	tracer := AuthErrorTracer + ": Register()"

	var payload dtos.ReqRegister
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(ce.ErrCodeInvalidPayload, ce.ErrMsgInvalidPayload, tracer, err)
		ctx.Error(err)
		return
	}

	data := entities.NewAuth{
		Email:    payload.Email,
		Password: payload.Password,
	}

	request := entities.NewRequest{
		UserAgent: ctx.Request.UserAgent(),
		IPAddress: ctx.ClientIP(),
	}

	token, err := h.au.Register(ctx, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dtos.RespAuth{
		Token: token.AccessToken,
	}

	utils.SetSessionCookie(ctx, token.SessionToken)
	utils.SetResponse(ctx, "registered successfully", response, http.StatusCreated)
}
