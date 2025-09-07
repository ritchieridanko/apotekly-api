package entities

type NewOAuth struct {
	Provider   int16
	UID        string
	Email      string
	IsVerified bool
}
