package phishlet

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

type EvaluationContext struct {
	UserAgent   string
	Email       string
	IPAddress   string
	Hostname    string
	Path        string
	CustomParams map[string]string
}

type ConditionEvaluator struct {
	geoIP GeoIPProvider
}

type GeoIPProvider interface {
	GetCountry(ip string) (string, error)
	GetRegion(ip string) (string, error)
}

type MockGeoIPProvider struct{}

func (m *MockGeoIPProvider) GetCountry(ip string) (string, error) {
	return "US", nil
}

func (m *MockGeoIPProvider) GetRegion(ip string) (string, error) {
	return "California", nil
}

func NewConditionEvaluator() *ConditionEvaluator {
	return &ConditionEvaluator{
		geoIP: &MockGeoIPProvider{},
	}
}

type Condition struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Values  []string          `json:"values,omitempty"`
	Regex   *regexp.Regexp    `json:"-"`
	Actions []ConditionAction `json:"actions"`
}

type ConditionAction struct {
	Type     string `json:"type"`
	Value    string `json:"value,omitempty"`
	Template string `json:"template,omitempty"`
}

func (ce *ConditionEvaluator) EvaluateConditions(conditions []Condition, context *EvaluationContext) ([]*ConditionAction, error) {
	var actions []*ConditionAction
	
	for _, condition := range conditions {
		matched, err := ce.evaluateCondition(&condition, context)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate condition '%s': %v", condition.Name, err)
		}
		
		if matched {
			for _, action := range condition.Actions {
				actionCopy := action
				actions = append(actions, &actionCopy)
			}
		}
	}
	
	return actions, nil
}

func (ce *ConditionEvaluator) evaluateCondition(condition *Condition, context *EvaluationContext) (bool, error) {
	switch strings.ToLower(condition.Type) {
	case "email_domain":
		return ce.evaluateEmailDomain(condition, context)
	case "user_agent":
		return ce.evaluateUserAgent(condition, context)
	case "ip_geo":
		return ce.evaluateIPGeo(condition, context)
	case "custom":
		return ce.evaluateCustom(condition, context)
	case "hostname":
		return ce.evaluateHostname(condition, context)
	case "path":
		return ce.evaluatePath(condition, context)
	default:
		return false, fmt.Errorf("unknown condition type: %s", condition.Type)
	}
}

func (ce *ConditionEvaluator) evaluateEmailDomain(condition *Condition, context *EvaluationContext) (bool, error) {
	if context.Email == "" {
		return false, nil
	}
	
	parts := strings.Split(context.Email, "@")
	if len(parts) != 2 {
		return false, nil
	}
	
	domain := strings.ToLower(parts[1])
	
	for _, value := range condition.Values {
		if strings.ToLower(value) == domain {
			return true, nil
		}
	}
	
	if condition.Regex != nil {
		return condition.Regex.MatchString(domain), nil
	}
	
	return false, nil
}

func (ce *ConditionEvaluator) evaluateUserAgent(condition *Condition, context *EvaluationContext) (bool, error) {
	if context.UserAgent == "" {
		return false, nil
	}
	
	userAgent := context.UserAgent
	
	for _, value := range condition.Values {
		if strings.Contains(strings.ToLower(userAgent), strings.ToLower(value)) {
			return true, nil
		}
	}
	
	if condition.Regex != nil {
		return condition.Regex.MatchString(userAgent), nil
	}
	
	return false, nil
}

func (ce *ConditionEvaluator) evaluateIPGeo(condition *Condition, context *EvaluationContext) (bool, error) {
	if context.IPAddress == "" {
		return false, nil
	}
	
	ip := net.ParseIP(context.IPAddress)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", context.IPAddress)
	}
	
	country, err := ce.geoIP.GetCountry(context.IPAddress)
	if err != nil {
		return false, err
	}
	
	for _, value := range condition.Values {
		if strings.EqualFold(value, country) {
			return true, nil
		}
	}
	
	if condition.Regex != nil {
		return condition.Regex.MatchString(country), nil
	}
	
	return false, nil
}

func (ce *ConditionEvaluator) evaluateCustom(condition *Condition, context *EvaluationContext) (bool, error) {
	if context.CustomParams == nil {
		return false, nil
	}
	
	for key, expectedValue := range context.CustomParams {
		for _, value := range condition.Values {
			parts := strings.SplitN(value, "=", 2)
			if len(parts) == 2 && parts[0] == key && parts[1] == expectedValue {
				return true, nil
			}
		}
	}
	
	return false, nil
}

func (ce *ConditionEvaluator) evaluateHostname(condition *Condition, context *EvaluationContext) (bool, error) {
	if context.Hostname == "" {
		return false, nil
	}
	
	hostname := strings.ToLower(context.Hostname)
	
	for _, value := range condition.Values {
		if strings.ToLower(value) == hostname {
			return true, nil
		}
	}
	
	if condition.Regex != nil {
		return condition.Regex.MatchString(hostname), nil
	}
	
	return false, nil
}

func (ce *ConditionEvaluator) evaluatePath(condition *Condition, context *EvaluationContext) (bool, error) {
	if context.Path == "" {
		return false, nil
	}
	
	path := context.Path
	
	for _, value := range condition.Values {
		if value == path {
			return true, nil
		}
	}
	
	if condition.Regex != nil {
		return condition.Regex.MatchString(path), nil
	}
	
	return false, nil
}
