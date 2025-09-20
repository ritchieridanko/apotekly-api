package constants

const (
	RedisKeyPasswordResetAuth  string = "password-reset:auth"
	RedisKeyPasswordResetToken string = "password-reset:token"
	RedisKeyTotalFailedAuth    string = "total-failed-auth"
	RedisKeyVerificationAuth   string = "verification:auth"
	RedisKeyVerificationToken  string = "verification:token"
)

const (
	RedisDurationTotalFailedAuth int = 15 // minutes
)

const (
	RedisMaxTotalFailedAuth int = 5 // times
)
