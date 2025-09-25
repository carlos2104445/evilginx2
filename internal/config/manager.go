package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Manager struct {
	config map[string]interface{}
	mu     sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		config: make(map[string]interface{}),
	}
}

func (m *Manager) LoadFromEnv() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.config["https_port"] = m.getEnvInt("EVILGINX_HTTPS_PORT", 443)
	m.config["dns_port"] = m.getEnvInt("EVILGINX_DNS_PORT", 53)
	m.config["session_timeout"] = m.getEnvDuration("EVILGINX_SESSION_TIMEOUT", "30m")
	m.config["rate_limit"] = m.getEnvInt("EVILGINX_RATE_LIMIT", 100)
	m.config["max_sessions"] = m.getEnvInt("EVILGINX_MAX_SESSIONS", 1000)
	m.config["cert_cache_size"] = m.getEnvInt("EVILGINX_CERT_CACHE_SIZE", 100)
	m.config["enable_debug"] = m.getEnvBool("EVILGINX_DEBUG", false)
	m.config["log_level"] = m.getEnvString("EVILGINX_LOG_LEVEL", "info")
}

func (m *Manager) Get(key string) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.config[key]
}

func (m *Manager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.config[key] = value
}

func (m *Manager) GetInt(key string, defaultValue int) int {
	if val := m.Get(key); val != nil {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return defaultValue
}

func (m *Manager) GetString(key string, defaultValue string) string {
	if val := m.Get(key); val != nil {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func (m *Manager) GetBool(key string, defaultValue bool) bool {
	if val := m.Get(key); val != nil {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

func (m *Manager) GetDuration(key string, defaultValue time.Duration) time.Duration {
	if val := m.Get(key); val != nil {
		if durVal, ok := val.(time.Duration); ok {
			return durVal
		}
	}
	return defaultValue
}

func (m *Manager) getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func (m *Manager) getEnvString(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func (m *Manager) getEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func (m *Manager) getEnvDuration(key string, defaultValue string) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	
	if duration, err := time.ParseDuration(val); err == nil {
		return duration
	}
	
	if defaultDuration, err := time.ParseDuration(defaultValue); err == nil {
		return defaultDuration
	}
	
	return 30 * time.Minute
}

func (m *Manager) Validate() error {
	httpsPort := m.GetInt("https_port", 443)
	if httpsPort < 1 || httpsPort > 65535 {
		return fmt.Errorf("invalid HTTPS port: %d", httpsPort)
	}
	
	dnsPort := m.GetInt("dns_port", 53)
	if dnsPort < 1 || dnsPort > 65535 {
		return fmt.Errorf("invalid DNS port: %d", dnsPort)
	}
	
	rateLimit := m.GetInt("rate_limit", 100)
	if rateLimit < 1 || rateLimit > 10000 {
		return fmt.Errorf("invalid rate limit: %d", rateLimit)
	}
	
	maxSessions := m.GetInt("max_sessions", 1000)
	if maxSessions < 1 || maxSessions > 100000 {
		return fmt.Errorf("invalid max sessions: %d", maxSessions)
	}
	
	return nil
}

func (m *Manager) GetAll() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	result := make(map[string]interface{})
	for k, v := range m.config {
		result[k] = v
	}
	
	return result
}
