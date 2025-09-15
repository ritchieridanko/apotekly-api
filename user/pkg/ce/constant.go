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
	ErrCodeMultimedia     int = 50009
)

// error tracers
const (
	AuthMiddlewareTracer string = "[MIDDLEWARE/AUTH(USER)]"
	DBTXTracer           string = "[DBTX(USER)]"
	UserHandlerTracer    string = "[HANDLER/USER]"
	UserRepoTracer       string = "[REPO/USER]"
	UserUsecaseTracer    string = "[USECASE/USER]"
)

// external error messages
const (
	ErrMsgInternalServer       string = "internal server error"
	ErrMsgInvalidCredentials   string = "invalid credentials"
	ErrMsgInvalidDataBirthdate string = "invalid data: birthdate"
	ErrMsgInvalidDataSex       string = "invalid data: sex"
	ErrMsgInvalidParams        string = "invalid params"
	ErrMsgInvalidPayload       string = "invalid payload"
	ErrMsgUserAlreadyExists    string = "user already exists"
	ErrMsgUnauthenticated      string = "unauthenticated"
)

// internal error loggers
var (
	ErrDBNoChange           error = fmt.Errorf("no rows affected after query execution")
	ErrInvalidAudience      error = fmt.Errorf("service not included in the token's audience")
	ErrInvalidDataBirthdate error = fmt.Errorf("invalid data-birthdate")
	ErrInvalidDataSex       error = fmt.Errorf("invalid data-sex")
	ErrInvalidTokenFormat   error = fmt.Errorf("invalid token format")
	ErrTokenNotFound        error = fmt.Errorf("token not found")
	ErrUserAlreadyExists    error = fmt.Errorf("auth already has user")
)
