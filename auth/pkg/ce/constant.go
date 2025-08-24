package ce

import "fmt"

// internal error codes
const (
	ErrCodeInvalidPayload int = 40000
	ErrCodeInvalidParams  int = 40001
	ErrCodeConflict       int = 40900
	ErrCodeDBQuery        int = 50000
	ErrCodeDBNoChange     int = 50001
	ErrCodePassHashing    int = 50002
	ErrCodeToken          int = 50003
	ErrCodeInvalidType    int = 50004
)

// error tracers
const (
	AuthHandlerTracer    string = "[HANDLER/AUTH]"
	AuthRepoTracer       string = "[REPO/AUTH]"
	AuthUsecaseTracer    string = "[USECASE/AUTH]"
	SessionRepoTracer    string = "[REPO/SESSION]"
	SessionUsecaseTracer string = "[USECASE/SESSION]"
)

// external error messages
const (
	ErrMsgEmailRegistered string = "email already registered"
	ErrMsgInternalServer  string = "internal server error"
	ErrMsgInvalidPayload  string = "invalid payload"
)

var (
	ErrDBNoChange             error = fmt.Errorf("no rows affected")
	ErrEmailAlreadyRegistered error = fmt.Errorf("email already registered")
	ErrInvalidType            error = fmt.Errorf("invalid return type")
)
