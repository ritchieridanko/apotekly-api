package ce

import (
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type internalErrorCode string

// internal error codes (for logs/debugging)
const (
	CodeAuthAudienceNotFound internalErrorCode = "AUTH_AUDIENCE_NOT_FOUND_ERROR"
	CodeAuthTokenExpired     internalErrorCode = "AUTH_TOKEN_EXPIRED_ERROR"
	CodeAuthTokenMalformed   internalErrorCode = "AUTH_TOKEN_MALFORMED_ERROR"
	CodeAuthTokenParsing     internalErrorCode = "AUTH_TOKEN_PARSING_ERROR"
	CodeAuthUnauthenticated  internalErrorCode = "AUTH_UNAUTHENTICATED_ERROR"
	CodeContextValueNotFound internalErrorCode = "CONTEXT_VALUE_NOT_FOUND_ERROR"
	CodeDBDuplicateData      internalErrorCode = "DB_DUPLICATE_DATA_ERROR"
	CodeDBQueryExecution     internalErrorCode = "DB_QUERY_EXECUTION_ERROR"
	CodeDBTransaction        internalErrorCode = "DB_TRANSACTION_ERROR"
	CodeInvalidPayload       internalErrorCode = "INVALID_PAYLOAD_ERROR"
	CodeUserNotFound         internalErrorCode = "USER_NOT_FOUND_ERROR"
)

// external error messages (for end-users)
const (
	MsgInternalServer  string = "Internal server error."
	MsgInvalidPayload  string = "Invalid payload."
	MsgUnauthenticated string = "Unauthenticated."
)

// internal error logs
var (
	ErrDBAffectNoRows error = errors.New("query execution affected no rows")
	ErrDBQueryNoRows  error = sql.ErrNoRows
	ErrTokenExpired   error = jwt.ErrTokenExpired
	ErrTokenMalformed error = jwt.ErrTokenMalformed
)
