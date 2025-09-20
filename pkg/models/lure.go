package models

import (
	"time"
)

type Lure struct {
	ID            string    `json:"id" db:"id"`
	Index         int       `json:"index" db:"index"`
	PhishletName  string    `json:"phishlet_name" db:"phishlet_name"`
	Hostname      string    `json:"hostname" db:"hostname"`
	Path          string    `json:"path" db:"path"`
	RedirectURL   string    `json:"redirect_url" db:"redirect_url"`
	OgTitle       string    `json:"og_title" db:"og_title"`
	OgDescription string    `json:"og_description" db:"og_description"`
	OgImageUrl    string    `json:"og_image_url" db:"og_image_url"`
	OgUrl         string    `json:"og_url" db:"og_url"`
	IsEnabled     bool      `json:"is_enabled" db:"is_enabled"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
	UpdateTime    time.Time `json:"update_time" db:"update_time"`
}

type LureStats struct {
	TotalLures   int `json:"total_lures"`
	EnabledLures int `json:"enabled_lures"`
	TotalClicks  int `json:"total_clicks"`
}
