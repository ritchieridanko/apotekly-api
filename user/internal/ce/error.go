package ce

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Error struct {
	Code    internalErrorCode
	Message string
	Err     error
}

func NewError(span trace.Span, code internalErrorCode, message string, err error) (newErr *Error) {
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

func (e *Error) Error() (err string) {
	return fmt.Sprintf("[%s] %s: %v\n", e.Code, e.Message, e.Err)
}

func (e *Error) ToExternalErrorCode() (externalCode int) {
	switch e.Code {
	case CodeInvalidPayload:
		return http.StatusBadRequest
	case CodeAuthAudienceNotFound, CodeAuthTokenExpired, CodeAuthTokenMalformed, CodeAuthUnauthenticated:
		return http.StatusUnauthorized
	case CodeDBDuplicateData:
		return http.StatusConflict
	case CodeAuthTokenParsing, CodeContextValueNotFound, CodeDBQueryExecution, CodeDBTransaction:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
