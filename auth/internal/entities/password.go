package entities

type ResetPassword struct {
	Token       string
	NewPassword string
}

type UpdatePassword struct {
	OldPassword string
	NewPassword string
}
