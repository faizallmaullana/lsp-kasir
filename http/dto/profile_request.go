package dto

type CreateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Contact  string `json:"contact"`
	Address  string `json:"address"`
	ImageUrl string `json:"image_url"`
}

type UpdateProfileRequest struct {
	Name     *string `json:"name"`
	Contact  *string `json:"contact"`
	Address  *string `json:"address"`
	ImageUrl *string `json:"image_url"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
