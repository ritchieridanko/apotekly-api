package constants

type ctxKey string

const (
	CtxKeyAuthID     ctxKey = "auth-id"
	CtxKeyIsVerified ctxKey = "is-verified"
	CtxKeyRequestID  ctxKey = "request-id"
	CtxKeyRoleID     ctxKey = "role-id"
)
