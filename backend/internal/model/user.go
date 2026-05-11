package model

import "time"

type User struct {
	ID                   int        `json:"id" db:"id"`
	Name                 string     `json:"name" db:"name"`
	Email                string     `json:"email" db:"email"`
	Password             string     `json:"-" db:"password"`
	Role                 string     `json:"role" db:"role"`
	Avatar               *string    `json:"avatar,omitempty" db:"avatar"`
	Bio                  *string    `json:"bio,omitempty" db:"bio"`
	Position             *string    `json:"position,omitempty" db:"position"`
	Department           *string    `json:"department,omitempty" db:"department"`
	IsActive             bool       `json:"is_active" db:"is_active"`
	EmailVerified        bool       `json:"email_verified" db:"email_verified"`
	OAuthProvider        *string    `json:"-" db:"oauth_provider"`
	OAuthID              *string    `json:"-" db:"oauth_id"`
	VerificationCode     *string    `json:"-" db:"verification_code"`
	VerificationSentAt   *time.Time `json:"-" db:"verification_sent_at"`
	VerificationAttempts int        `json:"-" db:"verification_attempts"`
	LastLoginAt          *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	NeedVerification bool   `json:"need_verification"`
	UserID           int    `json:"user_id,omitempty"`
	Token            string `json:"token,omitempty"`
	User             *User  `json:"user,omitempty"`
}

type AuthResponse struct {
	Token   string `json:"token,omitempty"`
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResendCodeRequest struct {
	Email string `json:"email"`
}

type OAuthRequest struct {
	Code     string `json:"code"`
	Provider string `json:"provider"`
}

type OAuthURLResponse struct {
	URL string `json:"url"`
}

type UpdateUserRequest struct {
	Name       *string `json:"name,omitempty"`
	Email      *string `json:"email,omitempty"`
	Avatar     *string `json:"avatar,omitempty"`
	Bio        *string `json:"bio,omitempty"`
	Position   *string `json:"position,omitempty"`
	Department *string `json:"department,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type ValidationError struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}
