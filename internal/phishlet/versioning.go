package phishlet

import (
	"context"
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
	storageVersions, err := pr.storage.ListPhishletVersions(ctx, name)
	if err != nil {
		return []*PhishletVersion{}, nil
	}
	
	var versions []*PhishletVersion
	for _, sv := range storageVersions {
		versions = append(versions, &PhishletVersion{
			Version:     sv.Version,
			Author:      sv.Author,
			Description: sv.Description,
			CreatedAt:   sv.CreatedAt,
			Hash:        sv.Hash,
			Content:     sv.Content,
		})
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
	
	newVersion := &storage.PhishletVersion{
		Version:     version,
		Author:      phishlet.Author,
		Description: description,
		CreatedAt:   time.Now().UTC(),
		Hash:        generateHash(phishlet.Name + version),
		Content:     "",
	}
	
	return pr.storage.CreatePhishletVersion(ctx, phishlet.Name, newVersion)
}

func (pr *PhishletRepository) CreateFlowSession(ctx context.Context, session *FlowSession) error {
	session.CreatedAt = time.Now().UTC()
	session.UpdatedAt = time.Now().UTC()
	
	storageSession := &storage.FlowSession{
		ID:           session.ID,
		PhishletName: session.PhishletName,
		FlowName:     session.FlowName,
		CurrentStep:  session.CurrentStep,
		StepData:     session.StepData,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}
	
	return pr.storage.CreateFlowSession(ctx, storageSession)
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
	storageSession, err := pr.storage.GetFlowSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("flow session not found: %s", sessionID)
	}
	
	session := &FlowSession{
		ID:           storageSession.ID,
		PhishletName: storageSession.PhishletName,
		FlowName:     storageSession.FlowName,
		CurrentStep:  storageSession.CurrentStep,
		StepData:     storageSession.StepData,
		CreatedAt:    storageSession.CreatedAt,
		UpdatedAt:    storageSession.UpdatedAt,
	}
	
	return session, nil
}

func generateHash(input string) string {
	return fmt.Sprintf("%x", []byte(input))[:8]
}
