package constants

type ctxKey string
type ctxKeyTx struct{}

const (
	CtxKeyAuthID ctxKey = "auth-id"
)

var (
	CtxKeyTx ctxKeyTx = ctxKeyTx{}
)
