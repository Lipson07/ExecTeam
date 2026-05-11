// internal/service/interfaces.go
package service

import "backend/internal/model"

type UserService interface {
	GetMe(userID int) (*model.User, error)
	UpdateMe(userID int, req *model.UpdateUserRequest) (*model.User, error)
	GetAll() ([]model.User, error)
	SearchByName(query string) ([]model.User, error)
	GetByID(id int) (*model.User, error)
	Delete(id int) error
}

type AuthService interface {
	Register(req *model.RegisterRequest) (*model.AuthResponse, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	VerifyEmail(email string, code string) (*model.AuthResponse, error)
	ResendCode(email string) error
	ValidateToken(tokenString string) (int, error)
}

type OAuthService interface {
	GetOAuthURL(provider string) (string, error)
	HandleCallback(provider string, code string) (*model.AuthResponse, error)
}
