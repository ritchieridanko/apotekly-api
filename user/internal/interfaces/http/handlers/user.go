package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/dto"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/validator"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/utils"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "handler.user"

type UserHandler struct {
	uu        usecases.UserUsecase
	validator *validator.Validator

	reqBodyMaxSize    int64
	allowedImageTypes []string
}

func NewUserHandler(
	uu usecases.UserUsecase,
	validator *validator.Validator,

	reqBodyMaxSize int64,
	allowedImageTypes []string,
) *UserHandler {
	return &UserHandler{uu, validator, reqBodyMaxSize, allowedImageTypes}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "CreateUser")
	defer span.End()

	// limit request body size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, h.reqBodyMaxSize)

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to create user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	var payload dto.CreateUserRequest
	if err := ctx.ShouldBindWith(&payload, binding.FormMultipart); err != nil {
		wErr := fmt.Errorf("failed to create user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		wErr := fmt.Errorf("failed to create user: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, wErr))
		return
	}

	normalizedBirthdate := payload.Birthdate.UTC()
	data := entities.CreateUser{
		Name:      payload.Name,
		Bio:       payload.Bio,
		Sex:       utils.NormalizePtr(payload.Sex),
		Birthdate: &normalizedBirthdate,
		Phone:     utils.TrimSpacePtr(payload.Phone),
	}

	var image multipart.File
	if payload.Image != nil {
		file, err := payload.Image.Open()
		if err != nil {
			log.Println("WARNING -> failed to open and upload image:", err.Error())
		} else {
			defer file.Close()

			if payload.Image.Size > h.reqBodyMaxSize {
				wErr := fmt.Errorf("failed to create user: %w", errors.New("file size exceeds max size"))
				msg := fmt.Sprintf("File size must not exceed %d bytes", h.reqBodyMaxSize)
				ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, msg, wErr))
				return
			}
			if err := h.validateImageType(ctxWithTracer, file); err != nil {
				ctx.Error(err)
				return
			}

			image = file
		}
	}

	user, err := h.uu.CreateUser(ctxWithTracer, authID, &data, image)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.CreateUserResponse{
		Created: h.userToResponse(*user),
	}

	utils.SetResponse(ctx, "User created successfully", response, http.StatusCreated)
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "GetUser")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	user, err := h.uu.GetUser(ctxWithTracer, authID)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := h.userToResponse(*user)

	utils.SetResponse(ctx, "User retrieved successfully", response, http.StatusOK)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "UpdateUser")
	defer span.End()

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to update user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	var payload dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		wErr := fmt.Errorf("failed to update user: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}
	if externalErr, internalErr := h.validator.Validate(payload); internalErr != nil {
		wErr := fmt.Errorf("failed to update user: %w", internalErr)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, externalErr, wErr))
		return
	}

	normalizedBirthdate := payload.Birthdate.UTC()
	data := entities.UpdateUser{
		Name:      payload.Name,
		Bio:       payload.Bio,
		Sex:       utils.NormalizePtr(payload.Sex),
		Birthdate: &normalizedBirthdate,
		Phone:     utils.TrimSpacePtr(payload.Phone),
	}

	user, err := h.uu.UpdateUser(ctxWithTracer, authID, &data)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.UpdateUserResponse{
		Updated: h.userToResponse(*user),
	}

	utils.SetResponse(ctx, "User updated successfully", response, http.StatusOK)
}

func (h *UserHandler) ChangeProfilePicture(ctx *gin.Context) {
	ctxWithTracer, span := otel.Tracer(userErrorTracer).Start(ctx.Request.Context(), "ChangeProfilePicture")
	defer span.End()

	// limit request body size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, h.reqBodyMaxSize)

	authID, err := utils.ContextGetAuthID(ctxWithTracer)
	if err != nil {
		wErr := fmt.Errorf("failed to change profile picture: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeContextValueNotFound, ce.MsgInternalServer, wErr))
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		wErr := fmt.Errorf("failed to change profile picture: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, wErr))
		return
	}

	image, err := file.Open()
	if err != nil {
		wErr := fmt.Errorf("failed to change profile picture: %w", err)
		ctx.Error(ce.NewError(span, ce.CodeRequestFile, ce.MsgInternalServer, wErr))
		return
	}
	defer image.Close()

	if file.Size > h.reqBodyMaxSize {
		wErr := fmt.Errorf("failed to change profile picture: %w", errors.New("file size exceeds max size"))
		msg := fmt.Sprintf("File size must not exceed %d bytes", h.reqBodyMaxSize)
		ctx.Error(ce.NewError(span, ce.CodeInvalidPayload, msg, wErr))
		return
	}
	if err := h.validateImageType(ctxWithTracer, image); err != nil {
		ctx.Error(err)
		return
	}

	user, err := h.uu.ChangeProfilePicture(ctxWithTracer, authID, image)
	if err != nil {
		ctx.Error(err)
		return
	}

	response := dto.UpdateUserResponse{
		Updated: h.userToResponse(*user),
	}

	utils.SetResponse(ctx, "Profile picture changed successfully", response, http.StatusOK)
}

func (h *UserHandler) userToResponse(user entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:             user.ID,
		Name:           user.Name,
		Bio:            user.Bio,
		Sex:            user.Sex,
		Birthdate:      user.Birthdate,
		Phone:          user.Phone,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func (h *UserHandler) validateImageType(ctx context.Context, image multipart.File) error {
	_, span := otel.Tracer(userErrorTracer).Start(ctx, "validateImageType")
	defer span.End()

	buf, err := io.ReadAll(image)
	if err != nil {
		wErr := fmt.Errorf("failed to validate image type: %w", err)
		return ce.NewError(span, ce.CodeFileBuffer, ce.MsgInternalServer, wErr)
	}

	// detect content type from first 512 bytes
	fileType := http.DetectContentType(buf[:utils.Min(len(buf), 512)])

	for _, allowedType := range h.allowedImageTypes {
		if fileType == allowedType {
			return nil
		}
	}

	err = fmt.Errorf("invalid file type '%s'", fileType)
	return ce.NewError(span, ce.CodeInvalidPayload, ce.MsgInvalidPayload, err)
}
