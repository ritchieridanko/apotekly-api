package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
)

func (h *authHandler) RegisterAsPharmacy(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(authErrorTracer).Start(ctx.Request.Context(), "RegisterAsPharmacy")
	defer span.End()

	var payload dto.ReqRegister
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		err := ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
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

	token, err := h.au.RegisterAsPharmacy(ctxWithTracer, &data, &request)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.RespAuth{
		Token: token.AccessToken,
	}

	utils.CookieSetSession(ctx, token.SessionToken)
	utils.SetResponse(ctx, "Link sent. Please check your email!", response, http.StatusCreated)
}
