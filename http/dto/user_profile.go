package dto

type MeResponse struct {
	UserID   string           `json:"user_id"`
	Email    string           `json:"email"`
	Role     string           `json:"role"`
	Profiles []ProfileSummary `json:"profiles"`
}

type ProfileSummary struct {
	IdProfile string `json:"id_profile"`
	Name      string `json:"name"`
	Contact   string `json:"contact"`
	Address   string `json:"address"`
	ImageUrl  string `json:"image_url"`
}
