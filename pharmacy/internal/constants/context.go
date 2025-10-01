package constants

type ctxKey string
type ctxKeyTx struct{}

const (
	CtxKeyAuthID     ctxKey = "auth-id"
	CtxKeyRoleID     ctxKey = "role-id"
	CtxKeyIsVerified ctxKey = "is-verified"
)

var (
	CtxKeyTx ctxKeyTx = ctxKeyTx{}
)
