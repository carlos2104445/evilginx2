# Proxy/C&C Separation Architecture

This document describes the Phase 2 implementation of Evilginx2 modernization: separating the monolithic HttpProxy component into standalone Proxy and C&C (Command & Control) services.

## Overview

The original `core/http_proxy.go` (1990 lines) has been refactored into two independent services:

- **Proxy Service**: Handles HTTP/HTTPS traffic forwarding, URL rewriting, content modification
- **C&C Service**: Manages sessions, phishlets, configuration, IP whitelisting

## Architecture

```
┌─────────────────┐    gRPC     ┌─────────────────┐
│  Proxy Service  │◄───────────►│  C&C Service    │
│                 │             │                 │
│ - HTTP/HTTPS    │             │ - Sessions      │
│ - URL Rewriting │             │ - Phishlets     │
│ - Content Patch │             │ - Config        │
│ - TLS Handling  │             │ - REST API      │
└─────────────────┘             └─────────────────┘
```

## Services

### Proxy Service (`cmd/proxy`)

**Purpose**: Stateless HTTP/HTTPS proxy that forwards traffic and modifies content.

**Key Features**:
- HTTP/HTTPS traffic interception and forwarding
- Host replacement (phish ↔ original domains)
- Content patching and JavaScript injection
- TLS certificate handling
- Request/response filtering

**Communication**: Queries C&C service via gRPC for:
- Phishlet lookups by hostname
- Session validation
- IP whitelist checks

**Usage**:
```bash
./proxy -port 443 -control localhost:8082 -certs ./certs
```

### C&C Service (`cmd/control-server`)

**Purpose**: Centralized management of sessions, phishlets, and configuration.

**Key Features**:
- Session lifecycle management
- Phishlet storage and retrieval
- IP whitelist management
- Configuration management
- gRPC API for proxy communication
- REST API for external clients (Phase 1 integration)

**Usage**:
```bash
./control-server -grpc-port 8082 -api-port 8081 -db ./control.db
```

## gRPC Communication

The services communicate via gRPC using the `ProxyControlService` defined in `proto/proxy_service.proto`:

### Service Methods

- `GetPhishletByHost`: Retrieve phishlet configuration by hostname
- `ValidateSession`: Check if a session should be handled
- `CreateSession`: Create new session
- `UpdateSession`: Update session data (credentials, tokens)
- `IsWhitelistedIP`: Check IP whitelist status
- `GetSessionIdByIP`: Get session ID for IP address
- `WhitelistIP`: Add IP to whitelist

### Message Types

- `Phishlet`: Complete phishlet configuration
- `ProxyHost`: Host mapping configuration
- `Session`: Session data and metadata
- Request/Response pairs for each service method

## Backward Compatibility

### Legacy Mode

The original `main.go` continues to work by starting both services together:

```bash
./evilginx2  # Starts both proxy and C&C in single process
```

### Migration Path

1. **Current**: Monolithic `main.go` with embedded HttpProxy
2. **Phase 2**: Separated services with legacy bridge
3. **Future**: Full microservices deployment

## Configuration

### Proxy Service Config (`config/proxy.yaml`)

```yaml
proxy:
  port: 443
  control_addr: "localhost:8082"
  cert_path: "./certs"
  
tls:
  auto_cert: true
  cert_cache: "./certs"
```

### C&C Service Config (`config/control.yaml`)

```yaml
control:
  grpc_port: 8082
  api_port: 8081
  database: "./control.db"
  
storage:
  type: "buntdb"
  path: "./control.db"
```

## Deployment Scenarios

### Single Machine (Legacy Compatible)
```bash
# Option 1: Legacy mode
./evilginx2

# Option 2: Separate services
./control-server -grpc-port 8082 -api-port 8081 &
./proxy -port 443 -control localhost:8082
```

### Distributed Deployment
```bash
# C&C Server (central management)
./control-server -grpc-port 8082 -api-port 8081

# Proxy Servers (edge locations)
./proxy -port 443 -control control-server:8082
./proxy -port 8443 -control control-server:8082
```

## Benefits

### Scalability
- Independent scaling of proxy and C&C components
- Multiple proxy instances can connect to single C&C
- Distributed deployment for global coverage

### Maintainability
- Clear separation of concerns
- Smaller, focused codebases
- Independent testing and deployment

### Performance
- gRPC communication for high performance
- Stateless proxy enables horizontal scaling
- Centralized session management reduces complexity

## Implementation Details

### Extracted Methods

**From HttpProxy to Proxy Service**:
- `httpsWorker()` - HTTPS connection handling
- `replaceHostWithOriginal()` - Host replacement logic
- `replaceHostWithPhished()` - Reverse host replacement
- `patchUrls()` - URL patching in content
- `injectJavascriptIntoBody()` - JS injection
- `injectOgHeaders()` - Open Graph header injection
- `blockRequest()` - Request blocking
- `setProxy()` - Proxy configuration

**From HttpProxy to C&C Service**:
- `setSessionUsername()` - Session credential storage
- `setSessionPassword()` - Session credential storage
- `setSessionCustom()` - Custom session data
- `whitelistIP()` - IP whitelist management
- `isWhitelistedIP()` - IP whitelist checking
- `getSessionIdByIP()` - Session lookup by IP
- `getPhishletByOrigHost()` - Phishlet lookup
- `getPhishletByPhishHost()` - Phishlet lookup
- `handleSession()` - Session lifecycle management

### Data Flow

1. **Request Processing**:
   ```
   Client → Proxy → gRPC(ValidateSession) → C&C
                 ← gRPC(Response) ←
          → Target Server
   ```

2. **Session Creation**:
   ```
   Proxy → gRPC(CreateSession) → C&C → Storage
   ```

3. **Configuration Updates**:
   ```
   Admin → REST API → C&C → Storage
                    → gRPC Notification → Proxy
   ```

## Testing

### Unit Tests
```bash
go test ./internal/proxy/...
go test ./internal/control/...
```

### Integration Tests
```bash
# Start C&C service
./control-server -grpc-port 8082 -api-port 8081 &

# Start proxy service
./proxy -port 8443 -control localhost:8082 &

# Test communication
curl -k https://localhost:8443/
curl http://localhost:8081/api/v1/health
```

### Load Testing
```bash
# Test proxy performance
ab -n 1000 -c 10 https://localhost:8443/

# Test C&C API performance
ab -n 1000 -c 10 http://localhost:8081/api/v1/phishlets
```

## Monitoring

### Metrics
- gRPC request/response times
- Session creation/update rates
- Proxy traffic volume
- Error rates and types

### Logging
- Structured logging with correlation IDs
- Service-specific log levels
- Centralized log aggregation support

## Security Considerations

### gRPC Security
- TLS encryption for inter-service communication
- Service authentication and authorization
- Rate limiting and request validation

### Network Security
- Firewall rules for service ports
- VPN/private network for inter-service communication
- Certificate management and rotation

## Future Enhancements

### Phase 3: Advanced Phishlet Features
- Conditional phishing logic
- Multi-page flow support
- Dynamic content adaptation

### Phase 4: Stealth & Evasion
- Intelligent traffic filtering
- Dynamic domain fronting
- Randomization and obfuscation

### Phase 5: Modern UI & DevOps
- Web-based management interface
- Container orchestration
- CI/CD pipeline integration

## Troubleshooting

### Common Issues

**gRPC Connection Failed**:
```bash
# Check C&C service is running
netstat -tlnp | grep 8082

# Check firewall rules
iptables -L | grep 8082
```

**Certificate Issues**:
```bash
# Check certificate directory permissions
ls -la ./certs/

# Regenerate certificates
rm -rf ./certs/* && ./proxy -port 443 -control localhost:8082
```

**Session Sync Issues**:
```bash
# Check database connectivity
./control-server -db ./control.db -grpc-port 8082

# Verify session storage
curl http://localhost:8081/api/v1/sessions
```

## Migration Guide

### From Monolithic to Separated Services

1. **Backup existing data**:
   ```bash
   cp -r ./data ./data.backup
   cp evilginx2.db evilginx2.db.backup
   ```

2. **Update configuration**:
   ```bash
   # Create service configs
   mkdir -p config/
   cp config/proxy.yaml.example config/proxy.yaml
   cp config/control.yaml.example config/control.yaml
   ```

3. **Start services**:
   ```bash
   # Start C&C first
   ./control-server -grpc-port 8082 -api-port 8081 &
   
   # Then start proxy
   ./proxy -port 443 -control localhost:8082
   ```

4. **Verify functionality**:
   ```bash
   # Test API
   curl http://localhost:8081/api/v1/health
   
   # Test proxy
   curl -k https://localhost:8443/
   ```

## Support

For issues related to Proxy/C&C separation:
1. Check service logs for error messages
2. Verify gRPC connectivity between services
3. Ensure proper configuration files
4. Test with legacy mode for comparison

For general Evilginx2 support, refer to the main documentation.
