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


type SubPhishlet struct {
	Name       string            `json:"name" db:"name"`
	ParentName string            `json:"parent_name" db:"parent_name"`
	Params     map[string]string `json:"params" db:"params"`
}
