package service

import (
	"backend/internal/model"
	repository "backend/internal/repo"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo         repository.UserRepo
	jwtSecret    []byte
	emailService *EmailService
}

func NewAuthService(repo repository.UserRepo, jwtSecret string, emailService *EmailService) AuthService {
	return &authService{
		repo:         repo,
		jwtSecret:    []byte(jwtSecret),
		emailService: emailService,
	}
}

func (s *authService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email уже используется")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &model.User{Name: req.Name, Email: req.Email, Password: string(hashed)}
	id, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	code := s.generateCode()
	s.repo.SetVerificationCode(id, code)
	s.emailService.SendVerificationCode(req.Email, code)

	return &model.AuthResponse{Message: "Код отправлен", Code: code}, nil
}

func (s *authService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("неверный email или пароль")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errors.New("неверный email или пароль")
	}

	if !user.EmailVerified {
		code := s.generateCode()
		s.repo.SetVerificationCode(user.ID, code)
		s.emailService.SendVerificationCode(user.Email, code)
		return &model.LoginResponse{NeedVerification: true, UserID: user.ID}, nil
	}

	s.repo.UpdateLastLogin(user.ID)
	token, _ := s.generateToken(user.ID)
	return &model.LoginResponse{Token: token, User: user}, nil
}

func (s *authService) VerifyEmail(email string, code string) (*model.AuthResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("пользователь не найден")
	}

	if user.VerificationCode == nil || *user.VerificationCode != code {
		s.repo.IncrementVerificationAttempts(user.ID)
		return nil, errors.New("неверный код")
	}

	if user.VerificationSentAt != nil && time.Since(*user.VerificationSentAt) > 15*time.Minute {
		return nil, errors.New("код просрочен")
	}

	s.repo.MarkEmailVerified(user.ID)
	s.repo.UpdateLastLogin(user.ID)
	token, _ := s.generateToken(user.ID)

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) ResendCode(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return errors.New("пользователь не найден")
	}
	if user.EmailVerified {
		return errors.New("email уже подтверждён")
	}

	code := s.generateCode()
	s.repo.SetVerificationCode(user.ID, code)
	s.emailService.SendVerificationCode(user.Email, code)
	return nil
}

func (s *authService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	userID, _ := claims["user_id"].(float64)
	return int(userID), nil
}

func (s *authService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *authService) generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
