package constants

type ctxKey string

const (
	CtxKeyAuthID     ctxKey = "auth-id"
	CtxKeyRoleID     ctxKey = "role-id"
	CtxKeyIsVerified ctxKey = "is-verified"
)
