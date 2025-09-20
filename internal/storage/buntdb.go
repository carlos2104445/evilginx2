package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kgretzky/evilginx2/pkg/models"
	"github.com/tidwall/buntdb"
)

const (
	SessionTable  = "sessions"
	PhishletTable = "phishlets"
	ConfigTable   = "config"
	LureTable     = "lures"
)

type BuntDBStorage struct {
	db   *buntdb.DB
	path string
}

func NewBuntDBStorage(path string) (*BuntDBStorage, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open buntdb: %w", err)
	}

	storage := &BuntDBStorage{
		db:   db,
		path: path,
	}

	if err := storage.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	return storage, nil
}

func (s *BuntDBStorage) init() error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		tx.CreateIndex("sessions_id", SessionTable+":*", buntdb.IndexJSON("id"))
		tx.CreateIndex("sessions_index", SessionTable+":*", buntdb.IndexJSON("index"))
		tx.CreateIndex("sessions_phishlet", SessionTable+":*", buntdb.IndexJSON("phishlet_name"))
		tx.CreateIndex("phishlets_name", PhishletTable+":*", buntdb.IndexJSON("name"))
		tx.CreateIndex("lures_id", LureTable+":*", buntdb.IndexJSON("id"))
		return nil
	})
}

func (s *BuntDBStorage) CreateSession(ctx context.Context, session *models.Session) error {
	if session.CreateTime.IsZero() {
		session.CreateTime = time.Now().UTC()
	}
	session.UpdateTime = time.Now().UTC()

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(SessionTable, session.ID)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) GetSession(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	err := s.db.View(func(tx *buntdb.Tx) error {
		key := s.genKey(SessionTable, id)
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &session)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

func (s *BuntDBStorage) GetSessionByIndex(ctx context.Context, index int) (*models.Session, error) {
	var session models.Session
	found := false

	err := s.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendEqual("sessions_index", s.getPivot(map[string]int{"index": index}), func(key, val string) bool {
			if err := json.Unmarshal([]byte(val), &session); err == nil {
				found = true
			}
			return false
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get session by index: %w", err)
	}
	if !found {
		return nil, fmt.Errorf("session not found with index: %d", index)
	}

	return &session, nil
}

func (s *BuntDBStorage) ListSessions(ctx context.Context, filters *SessionFilters) ([]*models.Session, error) {
	var sessions []*models.Session
	
	err := s.db.View(func(tx *buntdb.Tx) error {
		count := 0
		return tx.Ascend("sessions_id", func(key, val string) bool {
			if filters != nil && filters.Offset > 0 && count < filters.Offset {
				count++
				return true
			}
			
			if filters != nil && filters.Limit > 0 && len(sessions) >= filters.Limit {
				return false
			}

			var session models.Session
			if err := json.Unmarshal([]byte(val), &session); err == nil {
				if filters == nil || s.matchesSessionFilters(&session, filters) {
					sessions = append(sessions, &session)
				}
			}
			count++
			return true
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	return sessions, nil
}

func (s *BuntDBStorage) UpdateSession(ctx context.Context, session *models.Session) error {
	session.UpdateTime = time.Now().UTC()
	
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(SessionTable, session.ID)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) DeleteSession(ctx context.Context, id string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(SessionTable, id)
		_, err := tx.Delete(key)
		return err
	})
}

func (s *BuntDBStorage) CreatePhishlet(ctx context.Context, phishlet *models.Phishlet) error {
	if phishlet.CreateTime.IsZero() {
		phishlet.CreateTime = time.Now().UTC()
	}
	phishlet.UpdateTime = time.Now().UTC()

	data, err := json.Marshal(phishlet)
	if err != nil {
		return fmt.Errorf("failed to marshal phishlet: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(PhishletTable, phishlet.Name)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) GetPhishlet(ctx context.Context, name string) (*models.Phishlet, error) {
	var phishlet models.Phishlet
	err := s.db.View(func(tx *buntdb.Tx) error {
		key := s.genKey(PhishletTable, name)
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &phishlet)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get phishlet: %w", err)
	}
	return &phishlet, nil
}

func (s *BuntDBStorage) ListPhishlets(ctx context.Context, filters *PhishletFilters) ([]*models.Phishlet, error) {
	var phishlets []*models.Phishlet
	
	err := s.db.View(func(tx *buntdb.Tx) error {
		count := 0
		return tx.Ascend("phishlets_name", func(key, val string) bool {
			if filters != nil && filters.Offset > 0 && count < filters.Offset {
				count++
				return true
			}
			
			if filters != nil && filters.Limit > 0 && len(phishlets) >= filters.Limit {
				return false
			}

			var phishlet models.Phishlet
			if err := json.Unmarshal([]byte(val), &phishlet); err == nil {
				if filters == nil || s.matchesPhishletFilters(&phishlet, filters) {
					phishlets = append(phishlets, &phishlet)
				}
			}
			count++
			return true
		})
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to list phishlets: %w", err)
	}
	
	return phishlets, nil
}

func (s *BuntDBStorage) UpdatePhishlet(ctx context.Context, phishlet *models.Phishlet) error {
	phishlet.UpdateTime = time.Now().UTC()

	data, err := json.Marshal(phishlet)
	if err != nil {
		return fmt.Errorf("failed to marshal phishlet: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(PhishletTable, phishlet.Name)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) DeletePhishlet(ctx context.Context, name string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(PhishletTable, name)
		_, err := tx.Delete(key)
		return err
	})
}

func (s *BuntDBStorage) GetConfig(ctx context.Context) (*models.Config, error) {
	var config models.Config
	err := s.db.View(func(tx *buntdb.Tx) error {
		key := s.genKey(ConfigTable, "main")
		val, err := tx.Get(key)
		if err == buntdb.ErrNotFound {
			config = models.Config{
				General: models.GeneralConfig{
					HttpsPort: 443,
					DnsPort:   53,
				},
				UpdateTime: time.Now().UTC(),
			}
			return nil
		}
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &config)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return &config, nil
}

func (s *BuntDBStorage) UpdateConfig(ctx context.Context, config *models.Config) error {
	config.UpdateTime = time.Now().UTC()
	
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(ConfigTable, "main")
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) CreateLure(ctx context.Context, lure *models.Lure) error {
	if lure.CreateTime.IsZero() {
		lure.CreateTime = time.Now().UTC()
	}
	lure.UpdateTime = time.Now().UTC()

	data, err := json.Marshal(lure)
	if err != nil {
		return fmt.Errorf("failed to marshal lure: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(LureTable, lure.ID)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) GetLure(ctx context.Context, id string) (*models.Lure, error) {
	var lure models.Lure
	err := s.db.View(func(tx *buntdb.Tx) error {
		key := s.genKey(LureTable, id)
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &lure)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get lure: %w", err)
	}
	return &lure, nil
}

func (s *BuntDBStorage) ListLures(ctx context.Context) ([]*models.Lure, error) {
	var lures []*models.Lure
	
	err := s.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("lures_id", func(key, val string) bool {
			var lure models.Lure
			if err := json.Unmarshal([]byte(val), &lure); err == nil {
				lures = append(lures, &lure)
			}
			return true
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list lures: %w", err)
	}
	return lures, nil
}

func (s *BuntDBStorage) UpdateLure(ctx context.Context, lure *models.Lure) error {
	lure.UpdateTime = time.Now().UTC()
	
	data, err := json.Marshal(lure)
	if err != nil {
		return fmt.Errorf("failed to marshal lure: %w", err)
	}

	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(LureTable, lure.ID)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) DeleteLure(ctx context.Context, id string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey(LureTable, id)
		_, err := tx.Delete(key)
		return err
	})
}

func (s *BuntDBStorage) CreatePhishletVersion(ctx context.Context, name string, version *PhishletVersion) error {
	data, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("failed to marshal phishlet version: %w", err)
	}
	
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey("phishlet_version", name+":"+version.Version)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
}

func (s *BuntDBStorage) ListPhishletVersions(ctx context.Context, name string) ([]*PhishletVersion, error) {
	var versions []*PhishletVersion
	
	err := s.db.View(func(tx *buntdb.Tx) error {
		prefix := "phishlet_version:" + name + ":"
		return tx.Ascend("", func(key, val string) bool {
			if !strings.HasPrefix(key, prefix) {
				return true
			}
			
			var version PhishletVersion
			if err := json.Unmarshal([]byte(val), &version); err == nil {
				versions = append(versions, &version)
			}
			return true
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list phishlet versions: %w", err)
	}
	return versions, nil
}

func (s *BuntDBStorage) GetPhishletVersion(ctx context.Context, name, version string) (*models.Phishlet, error) {
	var phishletVersion PhishletVersion
	err := s.db.View(func(tx *buntdb.Tx) error {
		key := s.genKey("phishlet_version", name+":"+version)
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &phishletVersion)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get phishlet version: %w", err)
	}
	
	phishlet := &models.Phishlet{
		Name:        name,
		Author:      phishletVersion.Author,
		Version:     version,
		RedirectURL: "",
		IsTemplate:  false,
		CreateTime:  phishletVersion.CreatedAt,
		UpdateTime:  phishletVersion.CreatedAt,
	}
	
	return phishlet, nil
}

func (s *BuntDBStorage) CreateFlowSession(ctx context.Context, session *FlowSession) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal flow session: %w", err)
	}
	
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey("flow_session", session.ID)
		_, _, err := tx.Set(key, string(data), nil)
		return err
	})
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
	var session FlowSession
	err := s.db.View(func(tx *buntdb.Tx) error {
		key := s.genKey("flow_session", sessionID)
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &session)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get flow session: %w", err)
	}
	
	return &session, nil
}

func (s *BuntDBStorage) DeleteFlowSession(ctx context.Context, sessionID string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		key := s.genKey("flow_session", sessionID)
		_, err := tx.Delete(key)
		return err
	})
}

func (s *BuntDBStorage) Close() error {
	return s.db.Close()
}

func (s *BuntDBStorage) Flush() error {
	return s.db.Shrink()
}

func (s *BuntDBStorage) genKey(table, id string) string {
	return table + ":" + id
}

func (s *BuntDBStorage) getPivot(t interface{}) string {
	pivot, _ := json.Marshal(t)
	return string(pivot)
}

func (s *BuntDBStorage) matchesSessionFilters(session *models.Session, filters *SessionFilters) bool {
	if filters.PhishletName != "" && session.PhishletName != filters.PhishletName {
		return false
	}
	if filters.Username != "" && session.Username != filters.Username {
		return false
	}
	if filters.StartTime != nil && session.CreateTime.Before(*filters.StartTime) {
		return false
	}
	if filters.EndTime != nil && session.CreateTime.After(*filters.EndTime) {
		return false
	}
	return true
}

func (s *BuntDBStorage) matchesPhishletFilters(phishlet *models.Phishlet, filters *PhishletFilters) bool {
	if filters.Name != "" && phishlet.Name != filters.Name {
		return false
	}
	if filters.Enabled != nil && phishlet.IsEnabled != *filters.Enabled {
		return false
	}
	return true
}
