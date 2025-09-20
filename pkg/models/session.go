package models

import (
	"time"
)

type Session struct {
	ID          string            `json:"id" db:"id"`
	PhishletID  string            `json:"phishlet_id" db:"phishlet_id"`
	LureID      string            `json:"lure_id" db:"lure_id"`
	Username    string            `json:"username" db:"username"`
	Password    string            `json:"password" db:"password"`
	Custom      map[string]string `json:"custom" db:"custom"`
	Tokens      map[string]string `json:"tokens" db:"tokens"`
	RemoteAddr  string            `json:"remote_addr" db:"remote_addr"`
	UserAgent   string            `json:"user_agent" db:"user_agent"`
	IsDone      bool              `json:"is_done" db:"is_done"`
	IsAuthUrl   bool              `json:"is_auth_url" db:"is_auth_url"`
	IsLanding   bool              `json:"is_landing" db:"is_landing"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}
