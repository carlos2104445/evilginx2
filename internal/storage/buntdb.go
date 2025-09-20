package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/buntdb"
	"github.com/kgretzky/evilginx2/pkg/models"
)

type BuntDBStorage struct {
	db *buntdb.DB
}

func NewBuntDBStorage(path string) (*BuntDBStorage, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open BuntDB: %v", err)
	}
	
	return &BuntDBStorage{db: db}, nil
}

func (s *BuntDBStorage) Close() error {
	return s.db.Close()
}

func (s *BuntDBStorage) Set(ctx context.Context, key, value string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
}

func (s *BuntDBStorage) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	return value, err
}

func (s *BuntDBStorage) Delete(ctx context.Context, key string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
}

func (s *BuntDBStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	err := s.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, prefix) {
				keys = append(keys, key)
			}
			return true
		})
	})
	return keys, err
}

func (s *BuntDBStorage) CreatePhishlet(ctx context.Context, phishlet *models.Phishlet) error {
	phishlet.CreatedAt = time.Now().UTC()
	phishlet.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(phishlet)
	if err != nil {
		return fmt.Errorf("failed to marshal phishlet: %v", err)
	}
	
	key := fmt.Sprintf("phishlet:%s", phishlet.Name)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) GetPhishlet(ctx context.Context, name string) (*models.Phishlet, error) {
	key := fmt.Sprintf("phishlet:%s", name)
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("phishlet not found: %s", name)
	}
	
	var phishlet models.Phishlet
	if err := json.Unmarshal([]byte(data), &phishlet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal phishlet: %v", err)
	}
	
	return &phishlet, nil
}

func (s *BuntDBStorage) UpdatePhishlet(ctx context.Context, phishlet *models.Phishlet) error {
	phishlet.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(phishlet)
	if err != nil {
		return fmt.Errorf("failed to marshal phishlet: %v", err)
	}
	
	key := fmt.Sprintf("phishlet:%s", phishlet.Name)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) DeletePhishlet(ctx context.Context, name string) error {
	key := fmt.Sprintf("phishlet:%s", name)
	return s.Delete(ctx, key)
}

func (s *BuntDBStorage) ListPhishlets(ctx context.Context) ([]*models.Phishlet, error) {
	keys, err := s.List(ctx, "phishlet:")
	if err != nil {
		return nil, err
	}
	
	var phishlets []*models.Phishlet
	for _, key := range keys {
		data, err := s.Get(ctx, key)
		if err != nil {
			continue
		}
		
		var phishlet models.Phishlet
		if err := json.Unmarshal([]byte(data), &phishlet); err != nil {
			continue
		}
		
		phishlets = append(phishlets, &phishlet)
	}
	
	return phishlets, nil
}

func (s *BuntDBStorage) CreateSession(ctx context.Context, session *models.Session) error {
	session.CreatedAt = time.Now().UTC()
	session.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}
	
	key := fmt.Sprintf("session:%s", session.ID)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) GetSession(ctx context.Context, id string) (*models.Session, error) {
	key := fmt.Sprintf("session:%s", id)
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("session not found: %s", id)
	}
	
	var session models.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %v", err)
	}
	
	return &session, nil
}

func (s *BuntDBStorage) UpdateSession(ctx context.Context, session *models.Session) error {
	session.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}
	
	key := fmt.Sprintf("session:%s", session.ID)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) DeleteSession(ctx context.Context, id string) error {
	key := fmt.Sprintf("session:%s", id)
	return s.Delete(ctx, key)
}

func (s *BuntDBStorage) ListSessions(ctx context.Context) ([]*models.Session, error) {
	keys, err := s.List(ctx, "session:")
	if err != nil {
		return nil, err
	}
	
	var sessions []*models.Session
	for _, key := range keys {
		data, err := s.Get(ctx, key)
		if err != nil {
			continue
		}
		
		var session models.Session
		if err := json.Unmarshal([]byte(data), &session); err != nil {
			continue
		}
		
		sessions = append(sessions, &session)
	}
	
	return sessions, nil
}

func (s *BuntDBStorage) CreateLure(ctx context.Context, lure *models.Lure) error {
	lure.CreatedAt = time.Now().UTC()
	lure.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(lure)
	if err != nil {
		return fmt.Errorf("failed to marshal lure: %v", err)
	}
	
	key := fmt.Sprintf("lure:%s", lure.ID)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) GetLure(ctx context.Context, id string) (*models.Lure, error) {
	key := fmt.Sprintf("lure:%s", id)
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("lure not found: %s", id)
	}
	
	var lure models.Lure
	if err := json.Unmarshal([]byte(data), &lure); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lure: %v", err)
	}
	
	return &lure, nil
}

func (s *BuntDBStorage) UpdateLure(ctx context.Context, lure *models.Lure) error {
	lure.UpdatedAt = time.Now().UTC()
	
	data, err := json.Marshal(lure)
	if err != nil {
		return fmt.Errorf("failed to marshal lure: %v", err)
	}
	
	key := fmt.Sprintf("lure:%s", lure.ID)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) DeleteLure(ctx context.Context, id string) error {
	key := fmt.Sprintf("lure:%s", id)
	return s.Delete(ctx, key)
}

func (s *BuntDBStorage) ListLures(ctx context.Context) ([]*models.Lure, error) {
	keys, err := s.List(ctx, "lure:")
	if err != nil {
		return nil, err
	}
	
	var lures []*models.Lure
	for _, key := range keys {
		data, err := s.Get(ctx, key)
		if err != nil {
			continue
		}
		
		var lure models.Lure
		if err := json.Unmarshal([]byte(data), &lure); err != nil {
			continue
		}
		
		lures = append(lures, &lure)
	}
	
	return lures, nil
}

func (s *BuntDBStorage) SetConfig(ctx context.Context, key, value string) error {
	configKey := fmt.Sprintf("config:%s", key)
	return s.Set(ctx, configKey, value)
}

func (s *BuntDBStorage) GetConfig(ctx context.Context, key string) (string, error) {
	configKey := fmt.Sprintf("config:%s", key)
	return s.Get(ctx, configKey)
}

func (s *BuntDBStorage) DeleteConfig(ctx context.Context, key string) error {
	configKey := fmt.Sprintf("config:%s", key)
	return s.Delete(ctx, configKey)
}

func (s *BuntDBStorage) ListConfig(ctx context.Context) (map[string]string, error) {
	keys, err := s.List(ctx, "config:")
	if err != nil {
		return nil, err
	}
	
	config := make(map[string]string)
	for _, key := range keys {
		value, err := s.Get(ctx, key)
		if err != nil {
			continue
		}
		
		configKey := strings.TrimPrefix(key, "config:")
		config[configKey] = value
	}
	
	return config, nil
}

func (s *BuntDBStorage) CreatePhishletVersion(ctx context.Context, name string, version *PhishletVersion) error {
	data, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("failed to marshal phishlet version: %v", err)
	}
	
	key := fmt.Sprintf("phishlet_version:%s:%s", name, version.Version)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) ListPhishletVersions(ctx context.Context, name string) ([]*PhishletVersion, error) {
	prefix := fmt.Sprintf("phishlet_version:%s:", name)
	keys, err := s.List(ctx, prefix)
	if err != nil {
		return nil, err
	}
	
	var versions []*PhishletVersion
	for _, key := range keys {
		data, err := s.Get(ctx, key)
		if err != nil {
			continue
		}
		
		var version PhishletVersion
		if err := json.Unmarshal([]byte(data), &version); err != nil {
			continue
		}
		
		versions = append(versions, &version)
	}
	
	return versions, nil
}

func (s *BuntDBStorage) GetPhishletVersion(ctx context.Context, name, version string) (*models.Phishlet, error) {
	key := fmt.Sprintf("phishlet_version:%s:%s", name, version)
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("phishlet version not found: %s@%s", name, version)
	}
	
	var phishletVersion PhishletVersion
	if err := json.Unmarshal([]byte(data), &phishletVersion); err != nil {
		return nil, fmt.Errorf("failed to unmarshal phishlet version: %v", err)
	}
	
	phishlet := &models.Phishlet{
		Name:        name,
		Author:      phishletVersion.Author,
		Version:     version,
		RedirectURL: "",
		IsTemplate:  false,
		CreatedAt:   phishletVersion.CreatedAt,
		UpdatedAt:   phishletVersion.CreatedAt,
	}
	
	return phishlet, nil
}

func (s *BuntDBStorage) CreateFlowSession(ctx context.Context, session *FlowSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal flow session: %v", err)
	}
	
	key := fmt.Sprintf("flow_session:%s", session.ID)
	return s.Set(ctx, key, string(data))
}

func (s *BuntDBStorage) UpdateFlowSession(ctx context.Context, sessionID string, step string, data map[string]string) error {
	session, err := s.GetFlowSession(ctx, sessionID)
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
	
	return s.CreateFlowSession(ctx, session)
}

func (s *BuntDBStorage) GetFlowSession(ctx context.Context, sessionID string) (*FlowSession, error) {
	key := fmt.Sprintf("flow_session:%s", sessionID)
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("flow session not found: %s", sessionID)
	}
	
	var session FlowSession
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal flow session: %v", err)
	}
	
	return &session, nil
}

func (s *BuntDBStorage) DeleteFlowSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("flow_session:%s", sessionID)
	return s.Delete(ctx, key)
}
