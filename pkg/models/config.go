package models

import (
	"time"
)

type Config struct {
	General     GeneralConfig     `json:"general" db:"general"`
	Proxy       ProxyConfig       `json:"proxy" db:"proxy"`
	Blacklist   BlacklistConfig   `json:"blacklist" db:"blacklist"`
	GoPhish     GoPhishConfig     `json:"gophish" db:"gophish"`
	Phishlets   []PhishletConfig  `json:"phishlets" db:"phishlets"`
	UpdateTime  time.Time         `json:"update_time" db:"update_time"`
}

type GeneralConfig struct {
	Domain       string `json:"domain" db:"domain"`
	ExternalIPv4 string `json:"external_ipv4" db:"external_ipv4"`
	BindIPv4     string `json:"bind_ipv4" db:"bind_ipv4"`
	UnauthURL    string `json:"unauth_url" db:"unauth_url"`
	HttpsPort    int    `json:"https_port" db:"https_port"`
	DnsPort      int    `json:"dns_port" db:"dns_port"`
	Autocert     bool   `json:"autocert" db:"autocert"`
}

type ProxyConfig struct {
	Type     string `json:"type" db:"type"`
	Address  string `json:"address" db:"address"`
	Port     int    `json:"port" db:"port"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Enabled  bool   `json:"enabled" db:"enabled"`
}

type BlacklistConfig struct {
	Mode string `json:"mode" db:"mode"`
}

type GoPhishConfig struct {
	AdminURL    string `json:"admin_url" db:"admin_url"`
	ApiKey      string `json:"api_key" db:"api_key"`
	InsecureTLS bool   `json:"insecure_tls" db:"insecure_tls"`
}

type PhishletConfig struct {
	Name      string `json:"name" db:"name"`
	Hostname  string `json:"hostname" db:"hostname"`
	UnauthURL string `json:"unauth_url" db:"unauth_url"`
	Enabled   bool   `json:"enabled" db:"enabled"`
	Visible   bool   `json:"visible" db:"visible"`
}

type Lure struct {
	ID              string `json:"id" db:"id"`
	Hostname        string `json:"hostname" db:"hostname"`
	Path            string `json:"path" db:"path"`
	RedirectURL     string `json:"redirect_url" db:"redirect_url"`
	PhishletName    string `json:"phishlet_name" db:"phishlet_name"`
	Redirector      string `json:"redirector" db:"redirector"`
	UserAgentFilter string `json:"user_agent_filter" db:"user_agent_filter"`
	Info            string `json:"info" db:"info"`
	OgTitle         string `json:"og_title" db:"og_title"`
	OgDescription   string `json:"og_description" db:"og_description"`
	OgImageURL      string `json:"og_image_url" db:"og_image_url"`
	OgURL           string `json:"og_url" db:"og_url"`
	PausedUntil     int64  `json:"paused_until" db:"paused_until"`
	CreateTime      time.Time `json:"create_time" db:"create_time"`
	UpdateTime      time.Time `json:"update_time" db:"update_time"`
}

type SubPhishlet struct {
	Name       string            `json:"name" db:"name"`
	ParentName string            `json:"parent_name" db:"parent_name"`
	Params     map[string]string `json:"params" db:"params"`
}
