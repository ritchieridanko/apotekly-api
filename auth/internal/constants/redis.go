package constants

const (
	RedisKeyTotalFailedAuth    string = "total-failed-auth"
	RedisKeyPasswordResetAuth  string = "password-reset:auth"
	RedisKeyPasswordResetToken string = "password-reset:token"
)

const (
	RedisDurationTotalFailedAuth int = 15
)
