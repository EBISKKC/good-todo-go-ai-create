package input

type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

type LoginInput struct {
	Email    string
	Password string
}

type VerifyEmailInput struct {
	Token string
}

type RefreshTokenInput struct {
	RefreshToken string
}
