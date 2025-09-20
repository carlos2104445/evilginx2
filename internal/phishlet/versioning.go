package phishlet

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

type PhishletRepository struct {
	storage storage.Interface
	gitRepo string
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

func NewPhishletRepository(storage storage.Interface, gitRepo string) *PhishletRepository {
	return &PhishletRepository{
		storage: storage,
		gitRepo: gitRepo,
	}
}

func (pr *PhishletRepository) ListVersions(ctx context.Context, name string) ([]*PhishletVersion, error) {
	key := fmt.Sprintf("phishlet_versions:%s", name)
	data, err := pr.storage.Get(ctx, key)
	if err != nil {
		return []*PhishletVersion{}, nil
	}
	
	var versions []*PhishletVersion
	if err := json.Unmarshal([]byte(data), &versions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal versions: %v", err)
	}
	
	return versions, nil
}

func (pr *PhishletRepository) GetVersion(ctx context.Context, name, version string) (*models.Phishlet, error) {
	versions, err := pr.ListVersions(ctx, name)
	if err != nil {
		return nil, err
	}
	
	for _, v := range versions {
		if v.Version == version {
			phishlet := &models.Phishlet{
				Name:        name,
				Author:      v.Author,
				Version:     version,
				RedirectURL: "",
				IsTemplate:  false,
			}
			return phishlet, nil
		}
	}
	
	return nil, fmt.Errorf("version not found: %s", version)
}

func (pr *PhishletRepository) PublishVersion(ctx context.Context, phishlet *models.Phishlet, version string, description string) error {
	versions, err := pr.ListVersions(ctx, phishlet.Name)
	if err != nil {
		return err
	}
	
	for _, v := range versions {
		if v.Version == version {
			return fmt.Errorf("version already exists: %s", version)
		}
	}
	
	newVersion := &PhishletVersion{
		Version:     version,
		Author:      phishlet.Author,
		Description: description,
		CreatedAt:   time.Now().UTC(),
		Hash:        generateHash(phishlet.Name + version),
		Content:     "",
	}
	
	versions = append(versions, newVersion)
	
	data, err := json.Marshal(versions)
	if err != nil {
		return fmt.Errorf("failed to marshal versions: %v", err)
	}
	
	key := fmt.Sprintf("phishlet_versions:%s", phishlet.Name)
	return pr.storage.Set(ctx, key, string(data))
}

func (pr *PhishletRepository) CreateFlowSession(ctx context.Context, session *FlowSession) error {
	session.CreatedAt = time.Now().UTC()
	session.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal flow session: %v", err)
	}
	
	key := fmt.Sprintf("flow_session:%s", session.ID)
	return pr.storage.Set(ctx, key, string(data))
}

func (pr *PhishletRepository) UpdateFlowSession(ctx context.Context, sessionID string, step string, data map[string]string) error {
	session, err := pr.GetFlowSession(ctx, sessionID)
	if err != nil {
		return err
	}
	
	session.CurrentStep = step
	session.UpdatedAt = time.Now().UTC()
	
	if session.StepData == nil {
		session.StepData = make(map[string]string)
	}
	
	for k, v := range data {
		session.StepData[k] = v
	}
	
	return pr.CreateFlowSession(ctx, session)
}

func (pr *PhishletRepository) GetFlowSession(ctx context.Context, sessionID string) (*FlowSession, error) {
	key := fmt.Sprintf("flow_session:%s", sessionID)
	data, err := pr.storage.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("flow session not found: %s", sessionID)
	}
	
	var session FlowSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal flow session: %v", err)
	}
	
	return &session, nil
}

func generateHash(input string) string {
	return fmt.Sprintf("%x", []byte(input))[:8]
}
