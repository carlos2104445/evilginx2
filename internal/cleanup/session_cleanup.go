package cleanup

import (
	"context"
	"sync"
	"time"

	"github.com/kgretzky/evilginx2/log"
)

type SessionCleanup struct {
	sessions       map[string]*SessionInfo
	mu             sync.RWMutex
	cleanupTicker  *time.Ticker
	stopChan       chan struct{}
	sessionTimeout time.Duration
}

type SessionInfo struct {
	ID        string
	CreatedAt time.Time
	LastSeen  time.Time
	Active    bool
}

func NewSessionCleanup(sessionTimeout time.Duration) *SessionCleanup {
	return &SessionCleanup{
		sessions:       make(map[string]*SessionInfo),
		sessionTimeout: sessionTimeout,
		stopChan:       make(chan struct{}),
	}
}

func (sc *SessionCleanup) Start(ctx context.Context) {
	sc.cleanupTicker = time.NewTicker(5 * time.Minute)
	
	go func() {
		for {
			select {
			case <-sc.cleanupTicker.C:
				sc.cleanupExpiredSessions()
			case <-sc.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (sc *SessionCleanup) Stop() {
	if sc.cleanupTicker != nil {
		sc.cleanupTicker.Stop()
	}
	close(sc.stopChan)
}

func (sc *SessionCleanup) AddSession(sessionID string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	now := time.Now()
	sc.sessions[sessionID] = &SessionInfo{
		ID:        sessionID,
		CreatedAt: now,
		LastSeen:  now,
		Active:    true,
	}
}

func (sc *SessionCleanup) UpdateSession(sessionID string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	if session, exists := sc.sessions[sessionID]; exists {
		session.LastSeen = time.Now()
	}
}

func (sc *SessionCleanup) RemoveSession(sessionID string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	delete(sc.sessions, sessionID)
}

func (sc *SessionCleanup) cleanupExpiredSessions() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	now := time.Now()
	expiredSessions := []string{}
	
	for sessionID, session := range sc.sessions {
		if now.Sub(session.LastSeen) > sc.sessionTimeout {
			expiredSessions = append(expiredSessions, sessionID)
		}
	}
	
	for _, sessionID := range expiredSessions {
		delete(sc.sessions, sessionID)
		log.Info("Cleaned up expired session: %s", sessionID)
	}
	
	if len(expiredSessions) > 0 {
		log.Info("Cleaned up %d expired sessions", len(expiredSessions))
	}
}

func (sc *SessionCleanup) GetActiveSessionCount() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	
	return len(sc.sessions)
}

func (sc *SessionCleanup) GetSessionStats() map[string]interface{} {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	
	now := time.Now()
	activeCount := 0
	oldestSession := now
	
	for _, session := range sc.sessions {
		if session.Active {
			activeCount++
		}
		if session.CreatedAt.Before(oldestSession) {
			oldestSession = session.CreatedAt
		}
	}
	
	return map[string]interface{}{
		"total_sessions":  len(sc.sessions),
		"active_sessions": activeCount,
		"oldest_session":  oldestSession,
		"cleanup_timeout": sc.sessionTimeout.String(),
	}
}
