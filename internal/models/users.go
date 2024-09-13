package models

type CreateUserRequest struct {
	Email       string `json:"email"               validate:"required"`
	Password    string `json:"password"            validate:"required"`
}

type LoginUserRequest struct {
	Email       string `json:"email"               validate:"required"`
	Password    string `json:"password"            validate:"required"`
}