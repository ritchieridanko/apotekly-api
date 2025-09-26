package entities

type NewPassword struct {
	Token       string
	NewPassword string
}

type PasswordChange struct {
	OldPassword string
	NewPassword string
}
