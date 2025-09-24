package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var (
	hostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}
	
	if len(hostname) > 253 {
		return fmt.Errorf("hostname too long: %d characters (max 253)", len(hostname))
	}
	
	if !hostnameRegex.MatchString(hostname) {
		return fmt.Errorf("invalid hostname format: %s", hostname)
	}
	
	return nil
}

func ValidateIPAddress(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}
	
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address format: %s", ip)
	}
	
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	
	if len(email) > 254 {
		return fmt.Errorf("email too long: %d characters (max 254)", len(email))
	}
	
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	
	return nil
}

func SanitizeString(input string, maxLength int) string {
	if len(input) > maxLength {
		input = input[:maxLength]
	}
	
	input = strings.TrimSpace(input)
	
	input = strings.ReplaceAll(input, "\x00", "")
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, "\n", " ")
	
	return input
}

func ValidateRegexPattern(pattern string) error {
	if pattern == "" {
		return fmt.Errorf("regex pattern cannot be empty")
	}
	
	if len(pattern) > 1000 {
		return fmt.Errorf("regex pattern too long: %d characters (max 1000)", len(pattern))
	}
	
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}
	
	return nil
}
