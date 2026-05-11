package repository

import "backend/internal/model"

type UserRepo interface {
	Create(user *model.User) (int, error)
	FindByEmail(email string) (*model.User, error)
	FindByID(id int) (*model.User, error)
	FindByName(name string) ([]model.User, error)
	FindAll() ([]model.User, error)
	FindByOAuth(provider string, oauthID string) (*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
	UpdateLastLogin(id int) error
	SetVerificationCode(id int, code string) error
	FindByVerificationCode(id int, code string) (*model.User, error)
	MarkEmailVerified(id int) error
	IncrementVerificationAttempts(id int) error
}
