package models

import "time"

type Certificate struct {
	Domain     string    `json:"domain" db:"domain"`
	Issuer     string    `json:"issuer" db:"issuer"`
	NotBefore  time.Time `json:"not_before" db:"not_before"`
	NotAfter   time.Time `json:"not_after" db:"not_after"`
	IsValid    bool      `json:"is_valid" db:"is_valid"`
	IsWildcard bool      `json:"is_wildcard" db:"is_wildcard"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
