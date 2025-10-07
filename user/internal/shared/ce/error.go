package ce

import (
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Error struct {
	Code      errCode
	Message   string
	Err       error
	Timestamp time.Time
}

func NewError(span trace.Span, code errCode, message string, err error) *Error {
	if span != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, message)
	}

	return &Error{
		Code:      code,
		Message:   message,
		Err:       err,
		Timestamp: time.Now().UTC(),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s\t[%s]\t%v", e.Timestamp.Format("2006-01-02 15:04:05"), e.Code, e.Message)
}

func (e *Error) HTTPStatus() int {
	switch e.Code {
	case CodeInvalidParams, CodeInvalidPayload:
		return http.StatusBadRequest
	case
		CodeAuthAudienceNotFound,
		CodeAuthNotFound,
		CodeAuthTokenExpired,
		CodeAuthTokenMalformed,
		CodeAuthUnauthenticated,
		CodeInvalidTokenClaim,
		CodeRoleUnauthorized:
		return http.StatusUnauthorized
	case CodeAuthNotVerified:
		return http.StatusForbidden
	case CodeAddressNotFound, CodeUserNotFound:
		return http.StatusNotFound
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
