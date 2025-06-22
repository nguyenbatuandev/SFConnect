package entity

type RegsisterRequest struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role" validate:"required,oneof=admin partner buyer"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Message string `json:"message,omitempty" example:"Validation failed"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

