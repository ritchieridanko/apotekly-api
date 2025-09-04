package ce

import "fmt"

// internal error codes
const (
	ErrCodeInvalidPayload int = 40000
	ErrCodeInvalidParams  int = 40001
	ErrCodeInvalidAction  int = 40100
	ErrCodeConflict       int = 40900
	ErrCodeLocked         int = 42300
	ErrCodeDBQuery        int = 50000
	ErrCodeDBNoChange     int = 50001
	ErrCodePassHashing    int = 50002
	ErrCodeToken          int = 50003
	ErrCodeInvalidType    int = 50004
	ErrCodeCache          int = 50005
	ErrCodeDBTX           int = 50006
	ErrCodeContext        int = 50007
	ErrCodeParsing        int = 50008
	ErrCodeEmail          int = 50009
)

// error tracers
const (
	CacheTracer string = "[CACHE(AUTH)]"
	DBTXTracer  string = "[DBTX(AUTH)]"
	EmailTracer string = "[EMAIL(AUTH)]"

	AuthHandlerTracer    string = "[HANDLER/AUTH]"
	AuthMiddlewareTracer string = "[MIDDLEWARE/AUTH]"
	AuthRepoTracer       string = "[REPO/AUTH]"
	AuthUsecaseTracer    string = "[USECASE/AUTH]"
	OAuthRepoTracer      string = "[REPO/OAUTH]"
	SessionRepoTracer    string = "[REPO/SESSION]"
	SessionUsecaseTracer string = "[USECASE/SESSION]"
)

// external error messages
const (
	ErrMsgAccountLocked       string = "account locked"
	ErrMsgEmailRegistered     string = "email already registered"
	ErrMsgEmailNotVerified    string = "email not verified"
	ErrMsgInternalServer      string = "internal server error"
	ErrMsgInvalidCredentials  string = "invalid credentials"
	ErrMsgInvalidParams       string = "invalid params"
	ErrMsgInvalidPassword     string = "invalid password"
	ErrMsgInvalidPayload      string = "invalid payload"
	ErrMsgOAuthEmailChange    string = "oauth-linked accounts cannot change email"
	ErrMsgOAuthPasswordChange string = "oauth-linked accounts cannot change password"
	ErrMsgUnauthenticated     string = "unauthenticated"
)

// internal error loggers
var (
	ErrAccountLocked             error = fmt.Errorf("failed authentication multiple times, account locked")
	ErrDBNoChange                error = fmt.Errorf("no rows affected after query execution")
	ErrEmailAlreadyRegistered    error = fmt.Errorf("email already registered")
	ErrEmailVerificationRequired error = fmt.Errorf("email verification required")
	ErrInvalidTokenFormat        error = fmt.Errorf("invalid token format")
	ErrInvalidType               error = fmt.Errorf("invalid return type")
	ErrOAuthEmailChange          error = fmt.Errorf("oauth-linked account attempted to change email")
	ErrOAuthPasswordChange       error = fmt.Errorf("oauth-linked account attempted to change password")
	ErrOAuthRegularLogin         error = fmt.Errorf("regular auth attempted to authenticate by oauth")
	ErrSessionExpired            error = fmt.Errorf("session has expired")
	ErrSessionRevoked            error = fmt.Errorf("session has been revoked")
	ErrTokenEmpty                error = fmt.Errorf("token empty")
	ErrTokenNotFound             error = fmt.Errorf("token not found")
)
