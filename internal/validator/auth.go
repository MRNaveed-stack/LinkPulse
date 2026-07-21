package validator

type RegisterRequest struct {
	Username string `validate:"required,min=3,max=30"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=72"`
}
