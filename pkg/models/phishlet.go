package models

import (
	"time"
)

type Phishlet struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Author      string    `json:"author" db:"author"`
	Version     string    `json:"version" db:"version"`
	RedirectURL string    `json:"redirect_url" db:"redirect_url"`
	IsTemplate  bool      `json:"is_template" db:"is_template"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	
	Conditions     []Condition     `json:"conditions,omitempty"`
	MultiPageFlows []MultiPageFlow `json:"multi_page_flows,omitempty"`
}

type Condition struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"`
	Values  []string          `json:"values,omitempty"`
	Regex   string            `json:"regex,omitempty"`
	Actions []ConditionAction `json:"actions"`
}

type ConditionAction struct {
	Type     string `json:"type"`
	Value    string `json:"value,omitempty"`
	Template string `json:"template,omitempty"`
}

type MultiPageFlow struct {
	Name  string     `json:"name"`
	Steps []FlowStep `json:"steps"`
}

type FlowStep struct {
	Path        string            `json:"path"`
	Credentials []string          `json:"credentials"`
	NextStep    string            `json:"next_step,omitempty"`
	Conditions  map[string]string `json:"conditions,omitempty"`
}
