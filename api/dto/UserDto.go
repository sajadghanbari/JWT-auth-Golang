package dto

type CreateUpdateUserRequest struct {
	Name     string `json:"name" binding:"required,alpha,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
}

type UserResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}
