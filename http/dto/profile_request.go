package dto

// CreateProfileRequest represents payload to create a profile.
type CreateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Contact  string `json:"contact"`
	Address  string `json:"address"`
	ImageUrl string `json:"image_url"`
}

// UpdateProfileRequest represents payload to update profile fields.
type UpdateProfileRequest struct {
	Name     *string `json:"name"`
	Contact  *string `json:"contact"`
	Address  *string `json:"address"`
	ImageUrl *string `json:"image_url"`
}

// UpdateEmailRequest for updating user email.
type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
