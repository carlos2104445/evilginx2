package models

import (
	"time"
)

type Lure struct {
	ID          string    `json:"id" db:"id"`
	PhishletID  string    `json:"phishlet_id" db:"phishlet_id"`
	Hostname    string    `json:"hostname" db:"hostname"`
	Path        string    `json:"path" db:"path"`
	RedirectURL string    `json:"redirect_url" db:"redirect_url"`
	Info        string    `json:"info" db:"info"`
	OgTitle     string    `json:"og_title" db:"og_title"`
	OgDesc      string    `json:"og_desc" db:"og_desc"`
	OgImage     string    `json:"og_image" db:"og_image"`
	OgURL       string    `json:"og_url" db:"og_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
