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
	return fmt.Sprintf("[%s] -> %v", e.Code, e.Err)
}

func (e *Error) ToExternalErrorCode() (externalCode int) {
	switch e.Code {
	case CodeInvalidParams, CodeInvalidPayload:
		return http.StatusBadRequest
	case
		CodeAuthAudienceNotFound,
		CodeAuthTokenExpired,
		CodeAuthTokenMalformed,
		CodeAuthUnauthenticated:
		return http.StatusUnauthorized
	case CodeAuthNotVerified:
		return http.StatusForbidden
	case CodeDBDuplicateData:
		return http.StatusConflict
	case
		CodeAuthTokenParsing,
		CodeContextValueNotFound,
		CodeDBQueryExecution,
		CodeDBTransaction,
		CodeFileBuffer,
		CodeFileUploadFailed,
		CodeRequestFile:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
