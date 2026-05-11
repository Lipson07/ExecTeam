package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"backend/internal/model"
	repository "backend/internal/repo"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

type oAuthService struct {
	userRepo  repository.UserRepo
	jwtSecret []byte
	configs   map[string]OAuthConfig
}

type githubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type githubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type googleUser struct {
	ID      string `json:"sub"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

func NewOAuthService(userRepo repository.UserRepo, jwtSecret string, configs map[string]OAuthConfig) OAuthService {
	return &oAuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		configs:   configs,
	}
}

func (s *oAuthService) GetOAuthURL(provider string) (string, error) {
	cfg, ok := s.configs[provider]
	if !ok {
		return "", fmt.Errorf("неизвестный провайдер: %s", provider)
	}

	u, err := url.Parse(cfg.AuthURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURL)
	q.Set("response_type", "code")

	switch provider {
	case "github":
		q.Set("scope", "user:email")
	case "google":
		q.Set("scope", "openid email profile")
		q.Set("access_type", "offline")
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (s *oAuthService) HandleCallback(provider string, code string) (*model.AuthResponse, error) {
	cfg, ok := s.configs[provider]
	if !ok {
		return nil, fmt.Errorf("неизвестный провайдер: %s", provider)
	}

	accessToken, err := s.exchangeCode(cfg, code)
	if err != nil {
		return nil, err
	}

	email, name, avatar, oauthID, err := s.fetchUser(provider, cfg.UserInfoURL, accessToken)
	if err != nil {
		return nil, err
	}

	if email == "" {
		return nil, errors.New("не удалось получить email")
	}

	user, err := s.userRepo.FindByOAuth(provider, oauthID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, _ = s.userRepo.FindByEmail(email)
		if user != nil {
			user.OAuthProvider = &provider
			user.OAuthID = &oauthID
			if user.Avatar == nil && avatar != "" {
				user.Avatar = &avatar
			}
			s.userRepo.Update(user)
		} else {
			id, err := s.userRepo.Create(&model.User{
				Email:         email,
				Name:          name,
				OAuthProvider: &provider,
				OAuthID:       &oauthID,
			})
			if err != nil {
				return nil, err
			}
			s.userRepo.MarkEmailVerified(id)

			user = &model.User{
				ID:            id,
				Email:         email,
				Name:          name,
				OAuthProvider: &provider,
				OAuthID:       &oauthID,
			}
		}
	}

	if user.Avatar == nil && avatar != "" {
		user.Avatar = &avatar
	}

	s.userRepo.UpdateLastLogin(user.ID)

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}
func (s *oAuthService) exchangeCode(cfg OAuthConfig, code string) (string, error) {
	data := url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"redirect_uri":  {cfg.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.Post(cfg.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", errors.New("не удалось получить access_token")
	}

	return token, nil
}

func (s *oAuthService) fetchUser(provider string, url string, token string) (string, string, string, string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()

	switch provider {
	case "github":
		var u githubUser
		json.NewDecoder(resp.Body).Decode(&u)
		email := u.Email
		if email == "" {
			email, _ = s.fetchGitHubEmail(token)
		}
		name := u.Name
		if name == "" {
			name = u.Login
		}
		return email, name, u.AvatarURL, fmt.Sprintf("%d", u.ID), nil

	case "google":
		var u googleUser
		json.NewDecoder(resp.Body).Decode(&u)
		return u.Email, u.Name, u.Picture, u.ID, nil
	}

	return "", "", "", "", errors.New("неизвестный провайдер")
}

func (s *oAuthService) fetchGitHubEmail(token string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []githubEmail
	json.NewDecoder(resp.Body).Decode(&emails)

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", nil
}

func (s *oAuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
