# Evilginx2 Modernization Plan

## Current Architecture Analysis

### Monolithic Structure
The current evilginx2 codebase follows a monolithic architecture where all components are tightly coupled:

```
main.go
├── core/config.go (Configuration management)
├── core/http_proxy.go (Reverse proxy + session management)
├── core/phishlet.go (Phishlet engine)
├── core/terminal.go (CLI interface)
├── core/certdb.go (Certificate management)
├── core/nameserver.go (DNS server)
├── core/session.go (Session tracking)
├── core/blacklist.go (IP filtering)
└── database/ (Session storage)
```

### Key Dependencies & Coupling Issues

1. **HttpProxy** directly depends on:
   - Config, CertDb, Database, Blacklist, Session management
   - Phishlet loading and processing
   - Certificate generation and management

2. **Terminal** directly manages:
   - HttpProxy, Config, CertDb, Database
   - All CLI commands for phishlet/session management

3. **Phishlet** is tightly coupled to:
   - Config for parameter management
   - No versioning or repository system

## Proposed Modular Architecture

### Phase 1: API Foundation & Database Abstraction (Weeks 1-2)

#### 1.1 Create Core API Module
```
api/
├── server.go          # HTTP server setup
├── middleware.go      # Auth, CORS, logging
├── handlers/
│   ├── phishlets.go   # Phishlet CRUD operations
│   ├── sessions.go    # Session management
│   ├── config.go      # Configuration endpoints
│   └── certificates.go # Certificate management
└── models/
    ├── requests.go    # API request models
    └── responses.go   # API response models
```

#### 1.2 Database Abstraction Layer
```
storage/
├── interface.go       # Storage interface definition
├── buntdb/           # Current BuntDB implementation
│   └── sessions.go
├── postgres/         # Future PostgreSQL support
│   └── sessions.go
└── memory/           # In-memory for testing
    └── sessions.go
```

#### 1.3 Core API Endpoints
```
POST   /api/v1/phishlets                    # Create phishlet
GET    /api/v1/phishlets                    # List phishlets
GET    /api/v1/phishlets/{name}             # Get phishlet
PUT    /api/v1/phishlets/{name}             # Update phishlet
DELETE /api/v1/phishlets/{name}             # Delete phishlet
POST   /api/v1/phishlets/{name}/enable      # Enable phishlet
POST   /api/v1/phishlets/{name}/disable     # Disable phishlet

GET    /api/v1/sessions                     # List sessions
GET    /api/v1/sessions/{id}                # Get session details
DELETE /api/v1/sessions/{id}                # Delete session
POST   /api/v1/sessions/{id}/export         # Export session data

GET    /api/v1/config                       # Get configuration
PUT    /api/v1/config                       # Update configuration
GET    /api/v1/config/status                # System status

POST   /api/v1/certificates/generate        # Generate certificate
GET    /api/v1/certificates                 # List certificates
DELETE /api/v1/certificates/{domain}        # Delete certificate
```

### Phase 2: Proxy/C&C Separation (Weeks 3-4)

#### 2.1 Separate Proxy Service
```
services/
├── proxy/
│   ├── main.go           # Proxy service entry point
│   ├── server.go         # HTTP/HTTPS proxy server
│   ├── session.go        # Session handling
│   ├── phishlet.go       # Phishlet processing
│   └── middleware.go     # Request/response processing
└── control/
    ├── main.go           # C&C service entry point
    ├── api.go            # REST API server
    ├── terminal.go       # CLI interface
    └── management.go     # Proxy management
```

#### 2.2 Inter-Service Communication
```
communication/
├── grpc/
│   ├── proxy.proto       # Proxy service definitions
│   ├── control.proto     # Control service definitions
│   └── generated/        # Generated gRPC code
└── events/
    ├── publisher.go      # Event publishing
    ├── subscriber.go     # Event subscription
    └── types.go          # Event type definitions
```

#### 2.3 Configuration Management
```
config/
├── proxy.yaml           # Proxy service config
├── control.yaml         # Control service config
├── shared.yaml          # Shared configuration
└── loader.go            # Configuration loader
```

### Phase 3: Advanced Phishlet Features (Weeks 5-6)

#### 3.1 Enhanced Phishlet Engine
```
phishlet/
├── engine/
│   ├── processor.go      # Core phishlet processing
│   ├── conditional.go    # Conditional logic engine
│   ├── multipage.go      # Multi-page flow handler
│   └── validator.go      # Phishlet validation
├── repository/
│   ├── git.go           # Git-based repository
│   ├── local.go         # Local file system
│   └── remote.go        # Remote repository access
└── versioning/
    ├── manager.go       # Version management
    ├── updater.go       # Auto-update system
    └── compatibility.go # Version compatibility
```

#### 3.2 Conditional Phishing Logic
```yaml
# Enhanced phishlet format
conditional_flows:
  - condition:
      type: "email_domain"
      pattern: "@company\\.com$"
    actions:
      - redirect: "/corporate-login"
      - inject_js: "corporate-branding.js"
  - condition:
      type: "user_agent"
      pattern: "Mobile"
    actions:
      - template: "mobile-optimized"
```

#### 3.3 Multi-Page Flow Support
```yaml
multi_page_flows:
  - name: "microsoft_mfa"
    pages:
      - path: "/login"
        capture: ["username"]
        next_condition: "username_exists"
      - path: "/password"
        capture: ["password"]
        next_condition: "mfa_required"
      - path: "/mfa"
        capture: ["mfa_token"]
        completion: true
```

### Phase 4: Stealth & Evasion (Weeks 7-8)

#### 4.1 Intelligent Traffic Filtering
```
evasion/
├── filtering/
│   ├── bot_detector.go   # Bot detection algorithms
│   ├── fingerprint.go    # Browser fingerprinting
│   ├── reputation.go     # IP reputation checking
│   └── whitelist.go      # Whitelist management
├── obfuscation/
│   ├── url_randomizer.go # URL structure randomization
│   ├── payload_obfuscator.go # Payload obfuscation
│   └── timing_jitter.go  # Request timing variation
└── fronting/
    ├── domain_fronting.go # Domain fronting implementation
    ├── cdn_rotation.go    # CDN endpoint rotation
    └── proxy_chains.go    # Proxy chain management
```

#### 4.2 Advanced Bot Detection
```go
type BotDetector struct {
    UserAgentAnalyzer    *UserAgentAnalyzer
    BehaviorAnalyzer     *BehaviorAnalyzer
    FingerprintAnalyzer  *FingerprintAnalyzer
    ReputationChecker    *ReputationChecker
}

type DetectionResult struct {
    IsBot       bool
    Confidence  float64
    Reasons     []string
    Action      string // "block", "serve_decoy", "allow"
}
```

#### 4.3 Dynamic Obfuscation
```go
type ObfuscationEngine struct {
    URLRandomizer     *URLRandomizer
    PayloadObfuscator *PayloadObfuscator
    TimingJitter      *TimingJitter
}

// Randomize URL structures
func (e *ObfuscationEngine) RandomizeURL(baseURL string) string {
    // Add random query parameters, path segments
}

// Obfuscate JavaScript payloads
func (e *ObfuscationEngine) ObfuscatePayload(payload string) string {
    // Variable name randomization, code obfuscation
}
```

### Phase 5: Modern UI & DevOps (Weeks 9-10)

#### 5.1 Web UI Architecture
```
ui/
├── frontend/
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   └── utils/         # Utility functions
│   ├── public/
│   └── package.json
└── backend/
    └── static/            # Embedded static files
```

#### 5.2 Container Architecture
```
docker/
├── proxy/
│   └── Dockerfile        # Proxy service container
├── control/
│   └── Dockerfile        # Control service container
├── ui/
│   └── Dockerfile        # UI service container
└── docker-compose.yml    # Multi-service orchestration
```

#### 5.3 Deployment Configuration
```
deploy/
├── kubernetes/
│   ├── proxy-deployment.yaml
│   ├── control-deployment.yaml
│   ├── ui-deployment.yaml
│   └── ingress.yaml
├── helm/
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
└── terraform/
    ├── aws/              # AWS deployment
    ├── gcp/              # GCP deployment
    └── azure/            # Azure deployment
```

## Implementation Strategy

### Migration Approach

1. **Backward Compatibility**: Maintain existing CLI interface during transition
2. **Gradual Migration**: Implement new modules alongside existing code
3. **Feature Flags**: Use feature flags to enable/disable new functionality
4. **Testing Strategy**: Comprehensive testing for each phase

### File Structure After Modernization

```
evilginx2/
├── cmd/
│   ├── proxy/            # Proxy service binary
│   ├── control/          # Control service binary
│   └── evilginx/         # Legacy monolithic binary
├── internal/
│   ├── api/              # REST API implementation
│   ├── phishlet/         # Phishlet engine
│   ├── storage/          # Storage abstraction
│   ├── evasion/          # Stealth & evasion
│   └── communication/   # Inter-service communication
├── pkg/
│   ├── models/           # Shared data models
│   ├── config/           # Configuration management
│   └── utils/            # Utility functions
├── web/
│   ├── ui/               # Web interface
│   └── api/              # API documentation
├── deployments/
│   ├── docker/           # Container configurations
│   ├── kubernetes/       # K8s manifests
│   └── terraform/        # Infrastructure as code
├── phishlets/            # Phishlet repository
├── docs/                 # Documentation
└── scripts/              # Build and deployment scripts
```

## Benefits of Modularization

1. **Scalability**: Independent scaling of proxy and control services
2. **Maintainability**: Clear separation of concerns
3. **Extensibility**: Easy to add new features and integrations
4. **Testing**: Better unit and integration testing capabilities
5. **Deployment**: Flexible deployment options (monolithic or distributed)
6. **Security**: Improved security through service isolation
7. **Performance**: Optimized performance for specific use cases

## Risk Mitigation

1. **Complexity**: Start with simple modularization, add complexity gradually
2. **Performance**: Benchmark each phase to ensure no performance regression
3. **Compatibility**: Maintain backward compatibility throughout migration
4. **Testing**: Comprehensive testing strategy for each component
5. **Documentation**: Detailed documentation for new architecture

## Success Metrics

1. **Code Quality**: Reduced cyclomatic complexity, improved test coverage
2. **Performance**: No degradation in proxy performance
3. **Usability**: Improved user experience with web UI
4. **Maintainability**: Faster development cycles, easier bug fixes
5. **Scalability**: Ability to handle increased load through horizontal scaling
