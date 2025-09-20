package ce

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type internalErrorCode string

// internal error codes (for logs/debugging)
const (
	CodeAuthAudienceNotFound    internalErrorCode = "AUTH_AUDIENCE_NOT_FOUND_ERROR"
	CodeAuthEmailConflict       internalErrorCode = "AUTH_EMAIL_CONFLICT_ERROR"
	CodeAuthLocked              internalErrorCode = "AUTH_LOCKED_ERROR"
	CodeAuthNotFound            internalErrorCode = "AUTH_NOT_FOUND_ERROR"
	CodeAuthNotVerified         internalErrorCode = "AUTH_NOT_VERIFIED_ERROR"
	CodeAuthTokenExpired        internalErrorCode = "AUTH_TOKEN_EXPIRED_ERROR"
	CodeAuthTokenMalformed      internalErrorCode = "AUTH_TOKEN_MALFORMED_ERROR"
	CodeAuthTokenParsing        internalErrorCode = "AUTH_TOKEN_PARSING_ERROR"
	CodeAuthUnauthenticated     internalErrorCode = "AUTH_UNAUTHENTICATED_ERROR"
	CodeAuthVerified            internalErrorCode = "AUTH_VERIFIED_ERROR"
	CodeAuthWrongPassword       internalErrorCode = "AUTH_WRONG_PASSWORD_ERROR"
	CodeCacheBackoffWait        internalErrorCode = "CACHE_BACKOFF_WAIT_ERROR"
	CodeCacheQueryExecution     internalErrorCode = "CACHE_QUERY_EXECUTION_ERROR"
	CodeCacheValueNotFound      internalErrorCode = "CACHE_VALUE_NOT_FOUND_ERROR"
	CodeCacheScriptExecution    internalErrorCode = "CACHE_SCRIPT_EXECUTION_ERROR"
	CodeContextCookieNotFound   internalErrorCode = "CONTEXT_COOKIE_NOT_FOUND_ERROR"
	CodeContextValueNotFound    internalErrorCode = "CONTEXT_VALUE_NOT_FOUND_ERROR"
	CodeDBDuplicateData         internalErrorCode = "DB_DUPLICATE_DATA_ERROR"
	CodeDBQueryExecution        internalErrorCode = "DB_QUERY_EXECUTION_ERROR"
	CodeDBTransaction           internalErrorCode = "DB_TRANSACTION_ERROR"
	CodeEmailDelivery           internalErrorCode = "EMAIL_DELIVERY_ERROR"
	CodeEmailTemplateParsing    internalErrorCode = "EMAIL_TEMPLATE_PARSING_ERROR"
	CodeInvalidParams           internalErrorCode = "INVALID_PARAMS_ERROR"
	CodeInvalidPayload          internalErrorCode = "INVALID_PAYLOAD_ERROR"
	CodeJWTGenerationFailed     internalErrorCode = "JWT_GENERATION_FAILED_ERROR"
	CodeOAuthCodeExchangeFailed internalErrorCode = "OAUTH_CODE_EXCHANGE_FAILED_ERROR"
	CodeOAuthEmailChange        internalErrorCode = "OAUTH_EMAIL_CHANGE_ERROR"
	CodeOAuthGetUserInfoFailed  internalErrorCode = "OAUTH_GET_USER_INFO_FAILED_ERROR"
	CodeOAuthNotVerified        internalErrorCode = "OAUTH_NOT_VERIFIED_ERROR"
	CodeOAuthPasswordChange     internalErrorCode = "OAUTH_PASSWORD_CHANGE_ERROR"
	CodeOAuthRegularExists      internalErrorCode = "OAUTH_REGULAR_EXISTS_ERROR"
	CodeOAuthRegularLogin       internalErrorCode = "OAUTH_REGULAR_LOGIN_ERROR"
	CodePasswordHashingFailed   internalErrorCode = "PASSWORD_HASHING_FAILED_ERROR"
	CodeSessionExpired          internalErrorCode = "SESSION_EXPIRED_ERROR"
	CodeSessionNotFound         internalErrorCode = "SESSION_NOT_FOUND_ERROR"
	CodeSessionRevoked          internalErrorCode = "SESSION_REVOKED_ERROR"
	CodeTypeAssertionFailed     internalErrorCode = "TYPE_ASSERTION_FAILED_ERROR"
	CodeTypeConversionFailed    internalErrorCode = "TYPE_CONVERSION_FAILED_ERROR"
)

// external error messages (for end-users)
const (
	MsgInternalServer     string = "Internal server error."
	MsgInvalidCredentials string = "Invalid credentials."
	MsgInvalidParams      string = "Invalid params."
	MsgInvalidPayload     string = "Invalid payload."
	MsgInvalidToken       string = "Invalid token."
	MsgUnauthenticated    string = "Unauthenticated."
)

// internal error logs
var (
	ErrCookieNotFound       error = http.ErrNoCookie
	ErrDBAffectNoRows       error = errors.New("query execution affected no rows")
	ErrDBQueryNoRows        error = sql.ErrNoRows
	ErrSessionExpired       error = errors.New("session expired")
	ErrSessionRevoked       error = errors.New("session revoked")
	ErrTokenExpired         error = jwt.ErrTokenExpired
	ErrTokenMalformed       error = jwt.ErrTokenMalformed
	ErrTypeAssertionFailed  error = errors.New("type assertion failed")
	ErrTypeConversionFailed error = errors.New("type conversion failed")
)
