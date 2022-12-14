// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package client

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuthUser struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UserEmailStat struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	StatID    string    `json:"stat_id"`
	StatValue int32     `json:"stat_value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserEmailSyncHistory struct {
	UserID    uuid.UUID `json:"user_id"`
	HistoryID int64     `json:"history_id"`
	SyncedAt  time.Time `json:"synced_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserOauthToken struct {
	UserID    uuid.UUID       `json:"user_id"`
	Provider  string          `json:"provider"`
	Token     json.RawMessage `json:"token"`
	IsValid   bool            `json:"is_valid"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type UserProfile struct {
	UserID         uuid.UUID `json:"user_id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	IsActive       bool      `json:"is_active"`
	AutoArchive    bool      `json:"auto_archive"`
	AutoContribute bool      `json:"auto_contribute"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Waitlist struct {
	UserID           uuid.UUID       `json:"user_id"`
	Email            string          `json:"email"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	LinkedinUrl      string          `json:"linkedin_url"`
	Responses        json.RawMessage `json:"responses"`
	CanCreateAccount bool            `json:"can_create_account"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
