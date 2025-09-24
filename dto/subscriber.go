package dto

type CreateSubscriberRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateSubscriberRequest struct {
	Email string `json:"email" binding:"required,email"`
}
