package entities

type OAuth struct {
	Provider   int16
	UID        string
	Email      string
	IsVerified bool
}
