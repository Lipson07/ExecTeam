package model

import "time"

type GithubAccount struct {
	ID              int        `json:"id" db:"id"`
	UserID          int        `json:"user_id" db:"user_id"`
	GithubID        int        `json:"github_id" db:"github_id"`
	Username        string     `json:"username" db:"username"`
	ProfileURL      string     `json:"profile_url" db:"profile_url"`
	AvatarURL       string     `json:"avatar_url" db:"avatar_url"`
	PublicRepos     int        `json:"public_repos" db:"public_repos"`
	Followers       int        `json:"followers" db:"followers"`
	Following       int        `json:"following" db:"following"`
	Name            string     `json:"name" db:"name"`
	Company         string     `json:"company" db:"company"`
	Blog            string     `json:"blog" db:"blog"`
	Location        string     `json:"location" db:"location"`
	Bio             string     `json:"bio" db:"bio"`
	AccessToken     string     `json:"-" db:"access_token"`
	TokenExpiresAt  *time.Time `json:"token_expires_at,omitempty" db:"token_expires_at"`
	CreatedAtGithub *time.Time `json:"created_at_github,omitempty" db:"created_at_github"`
	UpdatedAtGithub *time.Time `json:"updated_at_github,omitempty" db:"updated_at_github"`
	AddedAt         time.Time  `json:"added_at" db:"added_at"`
	LastSyncedAt    *time.Time `json:"last_synced_at,omitempty" db:"last_synced_at"`
}
