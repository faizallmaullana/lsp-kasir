package dto

type CreateQuestionRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Question string `json:"question" binding:"required"`
}

type UpdateQuestionRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Question string `json:"question" binding:"required"`
}
