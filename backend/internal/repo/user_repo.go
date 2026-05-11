package repository

import (
	"database/sql"
	"fmt"

	"backend/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

var _ UserRepo = (*UserPostgres)(nil)

func (r *UserPostgres) Create(user *model.User) (int, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, role, is_active, created_at, updated_at
	`
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("UserPostgres.Create: %w", err)
	}
	return user.ID, nil
}

func (r *UserPostgres) FindByEmail(email string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserPostgres.FindByEmail: %w", err)
	}
	return &user, nil
}

func (r *UserPostgres) FindByID(id int) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserPostgres.FindByID: %w", err)
	}
	return &user, nil
}

func (r *UserPostgres) FindByName(name string) ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users WHERE name ILIKE '%' || $1 || '%'`
	err := r.db.Select(&users, query, name)
	if err != nil {
		return nil, fmt.Errorf("UserPostgres.FindByName: %w", err)
	}
	return users, nil
}

func (r *UserPostgres) FindAll() ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users ORDER BY id`
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("UserPostgres.FindAll: %w", err)
	}
	return users, nil
}

func (r *UserPostgres) FindByOAuth(provider string, oauthID string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE oauth_provider = $1 AND oauth_id = $2`
	err := r.db.Get(&user, query, provider, oauthID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserPostgres.FindByOAuth: %w", err)
	}
	return &user, nil
}

func (r *UserPostgres) FindOrCreateByOAuth(provider string, oauthID string, email string, name string) (*model.User, error) {
	user, err := r.FindByOAuth(provider, oauthID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	user, err = r.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		query := `UPDATE users SET oauth_provider = $1, oauth_id = $2 WHERE id = $3`
		_, err := r.db.Exec(query, provider, oauthID, user.ID)
		if err != nil {
			return nil, fmt.Errorf("UserPostgres.FindOrCreateByOAuth: %w", err)
		}
		return user, nil
	}

	newUser := &model.User{
		Name:          name,
		Email:         email,
		OAuthProvider: &provider,
		OAuthID:       &oauthID,
		EmailVerified: true,
	}

	id, err := r.Create(newUser)
	if err != nil {
		return nil, err
	}
	newUser.ID = id
	return newUser, nil
}

func (r *UserPostgres) Update(user *model.User) error {
	query := `
		UPDATE users 
		SET name = $1, email = $2, avatar = $3, bio = $4, position = $5, department = $6
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRow(query, user.Name, user.Email, user.Avatar, user.Bio, user.Position, user.Department, user.ID).
		Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("UserPostgres.Update: %w", err)
	}
	return nil
}

func (r *UserPostgres) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("UserPostgres.Delete: %w", err)
	}
	return nil
}

func (r *UserPostgres) UpdateLastLogin(id int) error {
	_, err := r.db.Exec(`UPDATE users SET last_login_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("UserPostgres.UpdateLastLogin: %w", err)
	}
	return nil
}

func (r *UserPostgres) SetVerificationCode(id int, code string) error {
	query := `UPDATE users SET verification_code = $1, verification_sent_at = NOW(), verification_attempts = 0 WHERE id = $2`
	_, err := r.db.Exec(query, code, id)
	if err != nil {
		return fmt.Errorf("UserPostgres.SetVerificationCode: %w", err)
	}
	return nil
}

func (r *UserPostgres) FindByVerificationCode(id int, code string) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id = $1 AND verification_code = $2`
	err := r.db.Get(&user, query, id, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("UserPostgres.FindByVerificationCode: %w", err)
	}
	return &user, nil
}

func (r *UserPostgres) MarkEmailVerified(id int) error {
	query := `UPDATE users SET email_verified = true, verification_code = NULL WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("UserPostgres.MarkEmailVerified: %w", err)
	}
	return nil
}

func (r *UserPostgres) IncrementVerificationAttempts(id int) error {
	query := `UPDATE users SET verification_attempts = verification_attempts + 1 WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("UserPostgres.IncrementVerificationAttempts: %w", err)
	}
	return nil
}
