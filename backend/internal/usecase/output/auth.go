package output

type RegisterOutput struct {
	UserID   string
	TenantID string
	Email    string
	Message  string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	TenantID     string
	Email        string
	Name         string
	Role         string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

type VerifyEmailOutput struct {
	Success bool
	Message string
}
