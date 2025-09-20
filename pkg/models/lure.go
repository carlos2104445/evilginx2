package models

import (
	"time"
)

type Lure struct {
	ID              string    `json:"id" db:"id"`
	Hostname        string    `json:"hostname" db:"hostname"`
	Path            string    `json:"path" db:"path"`
	RedirectURL     string    `json:"redirect_url" db:"redirect_url"`
	PhishletName    string    `json:"phishlet_name" db:"phishlet_name"`
	Redirector      string    `json:"redirector" db:"redirector"`
	UserAgentFilter string    `json:"user_agent_filter" db:"user_agent_filter"`
	Info            string    `json:"info" db:"info"`
	OgTitle         string    `json:"og_title" db:"og_title"`
	OgDescription   string    `json:"og_description" db:"og_description"`
	OgImageURL      string    `json:"og_image_url" db:"og_image_url"`
	OgURL           string    `json:"og_url" db:"og_url"`
	PausedUntil     int64     `json:"paused_until" db:"paused_until"`
	CreateTime      time.Time `json:"create_time" db:"create_time"`
	UpdateTime      time.Time `json:"update_time" db:"update_time"`
}

type LureStats struct {
	TotalLures   int `json:"total_lures"`
	EnabledLures int `json:"enabled_lures"`
	TotalClicks  int `json:"total_clicks"`
}
