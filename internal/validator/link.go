package validator

type CreateLinkRequest struct {
	Title          string `json:"title" validate:"required,min=1,max=100"`
	Slug           string `json:"slug" validate:"required,min=1,max=50,alphanum"`
	DestinationURL string `json:"destination_url" validate:"required,url"`
}
