package ce

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type errCode string

// internal error codes (for logs/debugging)
const (
	CodeAuthAudienceNotFound    errCode = "AUTH_AUDIENCE_NOT_FOUND_ERROR"
	CodeAuthEmailConflict       errCode = "AUTH_EMAIL_CONFLICT_ERROR"
	CodeAuthLocked              errCode = "AUTH_LOCKED_ERROR"
	CodeAuthNotFound            errCode = "AUTH_NOT_FOUND_ERROR"
	CodeAuthNotVerified         errCode = "AUTH_NOT_VERIFIED_ERROR"
	CodeAuthTokenExpired        errCode = "AUTH_TOKEN_EXPIRED_ERROR"
	CodeAuthTokenMalformed      errCode = "AUTH_TOKEN_MALFORMED_ERROR"
	CodeAuthTokenParsing        errCode = "AUTH_TOKEN_PARSING_ERROR"
	CodeAuthUnauthenticated     errCode = "AUTH_UNAUTHENTICATED_ERROR"
	CodeAuthVerified            errCode = "AUTH_VERIFIED_ERROR"
	CodeAuthWrongPassword       errCode = "AUTH_WRONG_PASSWORD_ERROR"
	CodeCacheBackoffWait        errCode = "CACHE_BACKOFF_WAIT_ERROR"
	CodeCacheQueryExecution     errCode = "CACHE_QUERY_EXECUTION_ERROR"
	CodeCacheValueNotFound      errCode = "CACHE_VALUE_NOT_FOUND_ERROR"
	CodeCacheScriptExecution    errCode = "CACHE_SCRIPT_EXECUTION_ERROR"
	CodeContextCookieNotFound   errCode = "CONTEXT_COOKIE_NOT_FOUND_ERROR"
	CodeContextValueNotFound    errCode = "CONTEXT_VALUE_NOT_FOUND_ERROR"
	CodeDBDuplicateData         errCode = "DB_DUPLICATE_DATA_ERROR"
	CodeDBQueryExecution        errCode = "DB_QUERY_EXECUTION_ERROR"
	CodeDBTransaction           errCode = "DB_TRANSACTION_ERROR"
	CodeEmailDelivery           errCode = "EMAIL_DELIVERY_ERROR"
	CodeEmailTemplateParsing    errCode = "EMAIL_TEMPLATE_PARSING_ERROR"
	CodeInvalidParams           errCode = "INVALID_PARAMS_ERROR"
	CodeInvalidPayload          errCode = "INVALID_PAYLOAD_ERROR"
	CodeInvalidTokenClaim       errCode = "INVALID_TOKEN_CLAIM_ERROR"
	CodeJWTGenerationFailed     errCode = "JWT_GENERATION_FAILED_ERROR"
	CodeOAuthCodeExchangeFailed errCode = "OAUTH_CODE_EXCHANGE_FAILED_ERROR"
	CodeOAuthEmailChange        errCode = "OAUTH_EMAIL_CHANGE_ERROR"
	CodeOAuthGetUserInfoFailed  errCode = "OAUTH_GET_USER_INFO_FAILED_ERROR"
	CodeOAuthNotVerified        errCode = "OAUTH_NOT_VERIFIED_ERROR"
	CodeOAuthPasswordChange     errCode = "OAUTH_PASSWORD_CHANGE_ERROR"
	CodeOAuthRegularExists      errCode = "OAUTH_REGULAR_EXISTS_ERROR"
	CodeOAuthRegularLogin       errCode = "OAUTH_REGULAR_LOGIN_ERROR"
	CodePasswordHashingFailed   errCode = "PASSWORD_HASHING_FAILED_ERROR"
	CodeRoleUnauthorized        errCode = "ROLE_UNAUTHORIZED_ERROR"
	CodeSessionExpired          errCode = "SESSION_EXPIRED_ERROR"
	CodeSessionNotFound         errCode = "SESSION_NOT_FOUND_ERROR"
	CodeSessionRevoked          errCode = "SESSION_REVOKED_ERROR"
	CodeTypeAssertionFailed     errCode = "TYPE_ASSERTION_FAILED_ERROR"
	CodeTypeConversionFailed    errCode = "TYPE_CONVERSION_FAILED_ERROR"
)

// external error messages (for end-users)
const (
	MsgEmailAlreadyRegistered string = "Email is already registered"
	MsgInternalServer         string = "Internal server error"
	MsgInvalidCredentials     string = "Invalid credentials"
	MsgInvalidParams          string = "Invalid params"
	MsgInvalidPayload         string = "Invalid payload"
	MsgInvalidToken           string = "Invalid token"
	MsgUnauthenticated        string = "Unauthenticated"
)

// internal error logs
var (
	ErrCacheNil            error = redis.Nil
	ErrCookieNotFound      error = http.ErrNoCookie
	ErrDBAffectNoRows      error = errors.New("query execution affected no rows")
	ErrDBQueryNoRows       error = sql.ErrNoRows
	ErrEmailConflict       error = errors.New("email conflict")
	ErrInvalidTokenClaim   error = errors.New("invalid token claim")
	ErrOAuthCodeNotFound   error = errors.New("oauth code not found")
	ErrSessionExpired      error = errors.New("session expired")
	ErrSessionRevoked      error = errors.New("session revoked")
	ErrTokenExpired        error = jwt.ErrTokenExpired
	ErrTokenMalformed      error = jwt.ErrTokenMalformed
	ErrTokenNotFound       error = errors.New("token not found")
	ErrTypeAssertionFailed error = errors.New("type assertion failed")
)
