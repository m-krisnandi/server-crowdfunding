package user

type RegisterUserInput struct {
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginInput struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}