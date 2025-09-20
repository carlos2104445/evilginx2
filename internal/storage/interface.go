package storage

import (
	"context"
	"time"

	"github.com/kgretzky/evilginx2/pkg/models"
)

type SessionFilters struct {
	PhishletName string
	Username     string
	StartTime    *time.Time
	EndTime      *time.Time
	Limit        int
	Offset       int
}

type PhishletFilters struct {
	Name    string
	Enabled *bool
	Limit   int
	Offset  int
}

type Interface interface {
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, id string) (*models.Session, error)
	GetSessionByIndex(ctx context.Context, index int) (*models.Session, error)
	ListSessions(ctx context.Context, filters *SessionFilters) ([]*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, id string) error

	CreatePhishlet(ctx context.Context, phishlet *models.Phishlet) error
	GetPhishlet(ctx context.Context, name string) (*models.Phishlet, error)
	ListPhishlets(ctx context.Context, filters *PhishletFilters) ([]*models.Phishlet, error)
	UpdatePhishlet(ctx context.Context, phishlet *models.Phishlet) error
	DeletePhishlet(ctx context.Context, name string) error

	GetConfig(ctx context.Context) (*models.Config, error)
	UpdateConfig(ctx context.Context, config *models.Config) error

	CreateLure(ctx context.Context, lure *models.Lure) error
	GetLure(ctx context.Context, id string) (*models.Lure, error)
	ListLures(ctx context.Context) ([]*models.Lure, error)
	UpdateLure(ctx context.Context, lure *models.Lure) error
	DeleteLure(ctx context.Context, id string) error

	CreatePhishletVersion(ctx context.Context, name string, version *PhishletVersion) error
	ListPhishletVersions(ctx context.Context, name string) ([]*PhishletVersion, error)
	GetPhishletVersion(ctx context.Context, name, version string) (*models.Phishlet, error)
	
	CreateFlowSession(ctx context.Context, session *FlowSession) error
	UpdateFlowSession(ctx context.Context, sessionID string, step string, data map[string]string) error
	GetFlowSession(ctx context.Context, sessionID string) (*FlowSession, error)
	DeleteFlowSession(ctx context.Context, sessionID string) error

	Close() error
	Flush() error
}

type PhishletVersion struct {
	Version     string    `json:"version"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Hash        string    `json:"hash"`
	Content     string    `json:"content"`
}

type FlowSession struct {
	ID           string            `json:"id"`
	PhishletName string            `json:"phishlet_name"`
	FlowName     string            `json:"flow_name"`
	CurrentStep  string            `json:"current_step"`
	StepData     map[string]string `json:"step_data"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
