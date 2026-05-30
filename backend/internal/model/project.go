package model

import "time"

type Project struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	GithubID    *int       `json:"github_id,omitempty" db:"github_id"`
	Status      string     `json:"status" db:"status"`
	OwnerID     int        `json:"owner_id" db:"owner_id"`
	StartDate   *time.Time `json:"start_date,omitempty" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Owner       *User      `json:"owner,omitempty" db:"-"`
	GithubRepo  *Github    `json:"github_repo,omitempty" db:"-"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	GithubID    *int   `json:"github_id,omitempty"`
	OwnerID     int    `json:"owner_id" validate:"required"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	GithubID    *int    `json:"github_id,omitempty"`
	Status      *string `json:"status,omitempty"`
	OwnerID     *int    `json:"owner_id,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}
