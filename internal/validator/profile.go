package validator

type UpdateProfileRequest struct {
	DisplayName string `validate:"required,max=100"`
	Bio         string `validate:"max=500"`
	AvatarURL   string `validate:"url"`
}
