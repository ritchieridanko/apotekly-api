package ce

import (
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type errCode string

// internal error codes (for logs/debugging)
const (
	CodeAddressNotFound      errCode = "ADDRESS_NOT_FOUND_ERROR"
	CodeAuthAudienceNotFound errCode = "AUTH_AUDIENCE_NOT_FOUND_ERROR"
	CodeAuthNotFound         errCode = "AUTH_NOT_FOUND_ERROR"
	CodeAuthNotVerified      errCode = "AUTH_NOT_VERIFIED_ERROR"
	CodeAuthTokenExpired     errCode = "AUTH_TOKEN_EXPIRED_ERROR"
	CodeAuthTokenMalformed   errCode = "AUTH_TOKEN_MALFORMED_ERROR"
	CodeAuthTokenParsing     errCode = "AUTH_TOKEN_PARSING_ERROR"
	CodeAuthUnauthenticated  errCode = "AUTH_UNAUTHENTICATED_ERROR"
	CodeContextValueNotFound errCode = "CONTEXT_VALUE_NOT_FOUND_ERROR"
	CodeDBDuplicateData      errCode = "DB_DUPLICATE_DATA_ERROR"
	CodeDBQueryExecution     errCode = "DB_QUERY_EXECUTION_ERROR"
	CodeDBTransaction        errCode = "DB_TRANSACTION_ERROR"
	CodeFileBuffer           errCode = "FILE_BUFFER_ERROR"
	CodeFileUploadFailed     errCode = "FILE_UPLOAD_FAILED_ERROR"
	CodeInvalidParams        errCode = "INVALID_PARAMS_ERROR"
	CodeInvalidPayload       errCode = "INVALID_PAYLOAD_ERROR"
	CodeInvalidTokenClaim    errCode = "INVALID_TOKEN_CLAIM_ERROR"
	CodeRequestFile          errCode = "REQUEST_FILE_ERROR"
	CodeRoleUnauthorized     errCode = "ROLE_UNAUTHORIZED_ERROR"
	CodeUserNotFound         errCode = "USER_NOT_FOUND_ERROR"
)

// external error messages (for end-users)
const (
	MsgAddressNotFound    string = "Address not found"
	MsgInternalServer     string = "Internal server error"
	MsgInvalidCredentials string = "Invalid credentials"
	MsgInvalidParams      string = "Invalid params"
	MsgInvalidPayload     string = "Invalid payload"
	MsgNoFieldsToUpdate   string = "No fields to update"
	MsgUnauthenticated    string = "Unauthenticated"
	MsgUserNotFound       string = "User not found"
)

// internal error logs
var (
	ErrDBAffectNoRows    error = errors.New("no rows affected")
	ErrDBQueryNoRows     error = sql.ErrNoRows
	ErrFileBufferRead    error = errors.New("buffer read failed")
	ErrInvalidFileType   error = errors.New("invalid file type")
	ErrInvalidTokenClaim error = errors.New("invalid token claim")
	ErrNoFieldsProvided  error = errors.New("no fields provided")
	ErrTokenExpired      error = jwt.ErrTokenExpired
	ErrTokenMalformed    error = jwt.ErrTokenMalformed
)
