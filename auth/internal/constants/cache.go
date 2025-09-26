package constants

const (
	CacheKeyEmailChangeAuth    string = "email-change:auth"
	CacheKeyEmailChangeToken   string = "email-change:token"
	CacheKeyPasswordResetAuth  string = "password-reset:auth"
	CacheKeyPasswordResetToken string = "password-reset:token"
	CacheKeyTotalFailedAuth    string = "total-failed-auth"
	CacheKeyVerificationAuth   string = "verification:auth"
	CacheKeyVerificationToken  string = "verification:token"
)

const (
	CacheDurationTotalFailedAuth int = 15 // minutes
)

const (
	CacheMaxTotalFailedAuth int = 5 // times
)
