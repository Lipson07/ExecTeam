package model

import "time"

type Github struct {
	ID              int        `json:"id" db:"id"`
	GithubAccountID int        `json:"github_account_id" db:"github_account_id"`
	RepoURL         string     `json:"repo_url" db:"repo_url"`
	RepoName        string     `json:"repo_name" db:"repo_name"`
	Owner           string     `json:"owner" db:"owner"`
	FullName        string     `json:"full_name" db:"full_name"`
	DefaultBranch   string     `json:"default_branch" db:"default_branch"`
	IsPrivate       bool       `json:"is_private" db:"is_private"`
	IsFork          bool       `json:"is_fork" db:"is_fork"`
	StarsCount      int        `json:"stars_count" db:"stars_count"`
	ForksCount      int        `json:"forks_count" db:"forks_count"`
	WatchersCount   int        `json:"watchers_count" db:"watchers_count"`
	OpenIssuesCount int        `json:"open_issues_count" db:"open_issues_count"`
	CommitsCount    int        `json:"commits_count" db:"commits_count"`
	PrimaryLanguage *string    `json:"primary_language,omitempty" db:"primary_language"`
	Languages       *string    `json:"languages,omitempty" db:"languages"`
	Description     *string    `json:"description,omitempty" db:"description"`
	Topics          *string    `json:"topics,omitempty" db:"topics"`
	License         *string    `json:"license,omitempty" db:"license"`
	HomepageURL     *string    `json:"homepage_url,omitempty" db:"homepage_url"`
	CloneURL        *string    `json:"clone_url,omitempty" db:"clone_url"`
	CreatedAtGithub *time.Time `json:"created_at_github,omitempty" db:"created_at_github"`
	UpdatedAtGithub *time.Time `json:"updated_at_github,omitempty" db:"updated_at_github"`
	LastPushAt      *time.Time `json:"last_push_at,omitempty" db:"last_push_at"`
	AddedAt         time.Time  `json:"added_at" db:"added_at"`
	LastSyncedAt    *time.Time `json:"last_synced_at,omitempty" db:"last_synced_at"`
	SizeKB          int        `json:"size_kb" db:"size_kb"`
	Archived        bool       `json:"archived" db:"archived"`
}
