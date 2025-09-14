package entities

type NewPassword struct {
	OldPassword string
	NewPassword string
}

type PasswordReset struct {
	Token       string
	NewPassword string
}
