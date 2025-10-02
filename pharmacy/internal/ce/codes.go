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
	CodeAuthNotFound         internalErrorCode = "AUTH_NOT_FOUND_ERROR"
	CodeAuthNotVerified      internalErrorCode = "AUTH_NOT_VERIFIED_ERROR"
	CodeAuthTokenExpired     internalErrorCode = "AUTH_TOKEN_EXPIRED_ERROR"
	CodeAuthTokenMalformed   internalErrorCode = "AUTH_TOKEN_MALFORMED_ERROR"
	CodeAuthTokenParsing     internalErrorCode = "AUTH_TOKEN_PARSING_ERROR"
	CodeAuthUnauthenticated  internalErrorCode = "AUTH_UNAUTHENTICATED_ERROR"
	CodeContextValueNotFound internalErrorCode = "CONTEXT_VALUE_NOT_FOUND_ERROR"
	CodeDBDuplicateData      internalErrorCode = "DB_DUPLICATE_DATA_ERROR"
	CodeDBQueryExecution     internalErrorCode = "DB_QUERY_EXECUTION_ERROR"
	CodeDBTransaction        internalErrorCode = "DB_TRANSACTION_ERROR"
	CodeFileBuffer           internalErrorCode = "FILE_BUFFER_ERROR"
	CodeFileUploadFailed     internalErrorCode = "FILE_UPLOAD_FAILED_ERROR"
	CodeInvalidParams        internalErrorCode = "INVALID_PARAMS_ERROR"
	CodeInvalidPayload       internalErrorCode = "INVALID_PAYLOAD_ERROR"
	CodePharmacyNotFound     internalErrorCode = "PHARMACY_NOT_FOUND_ERROR"
	CodeRequestFile          internalErrorCode = "REQUEST_FILE_ERROR"
	CodeRoleUnauthorized     internalErrorCode = "ROLE_UNAUTHORIZED_ERROR"
)

// external error messages (for end-users)
const (
	MsgInternalServer     string = "Internal server error."
	MsgInvalidCredentials string = "Invalid credentials."
	MsgInvalidParams      string = "Invalid params."
	MsgInvalidPayload     string = "Invalid payload."
	MsgNoFieldsToUpdate   string = "No fields to update."
	MsgPharmacyNotFound   string = "Pharmacy not found."
	MsgUnauthenticated    string = "Unauthenticated."
)

// internal error logs
var (
	ErrDBAffectNoRows   error = errors.New("query execution affected no rows")
	ErrDBQueryNoRows    error = sql.ErrNoRows
	ErrNoFieldsProvided error = errors.New("no fields provided")
	ErrTokenExpired     error = jwt.ErrTokenExpired
	ErrTokenMalformed   error = jwt.ErrTokenMalformed
)
