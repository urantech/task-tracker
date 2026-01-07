package user

import "time"

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	RoleId    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	PasswordRaw string `json:"password_raw" validate:"required,min=8"`
}

type RegisterResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}
