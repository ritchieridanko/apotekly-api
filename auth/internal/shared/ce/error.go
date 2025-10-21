package ce

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Error struct {
	Code    errCode
	Message string
	Err     error
}

func NewError(span trace.Span, code errCode, message string, err error) *Error {
	if span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
	}

	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) HTTPStatus() int {
	switch e.Code {
	case CodeAuthVerified, CodeCacheValueNotFound, CodeInvalidParams, CodeInvalidPayload:
		return http.StatusBadRequest
	case
		CodeAuthAudienceNotFound,
		CodeAuthNotFound,
		CodeAuthTokenExpired,
		CodeAuthTokenMalformed,
		CodeAuthUnauthenticated,
		CodeAuthWrongPassword,
		CodeContextCookieNotFound,
		CodeInvalidTokenClaim,
		CodeRoleUnauthorized,
		CodeSessionExpired,
		CodeSessionNotFound,
		CodeSessionRevoked:
		return http.StatusUnauthorized
	case
		CodeAuthNotVerified,
		CodeOAuthEmailChange,
		CodeOAuthNotVerified,
		CodeOAuthPasswordChange,
		CodeOAuthRegularLogin:
		return http.StatusForbidden
	case CodeAuthEmailConflict, CodeDBDuplicateData, CodeOAuthRegularExists:
		return http.StatusConflict
	case CodeAuthLocked:
		return http.StatusLocked
	case
		CodeAuthTokenParsing,
		CodeCacheBackoffWait,
		CodeCacheQueryExecution,
		CodeCacheScriptExecution,
		CodeContextValueNotFound,
		CodeDBQueryExecution,
		CodeDBTransaction,
		CodeEmailDelivery,
		CodeEmailTemplateParsing,
		CodeEventPublishingFailed,
		CodeJWTGenerationFailed,
		CodeOAuthCodeExchangeFailed,
		CodeOAuthGetUserInfoFailed,
		CodePasswordHashingFailed,
		CodeTypeAssertionFailed,
		CodeTypeConversionFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
