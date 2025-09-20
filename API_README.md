# Evilginx2 API Documentation

## Overview

The Evilginx2 Control API provides a RESTful interface for managing phishlets, sessions, configuration, and certificates. This API is part of Phase 1 of the modernization effort to create a modular, scalable architecture.

## Base URL

```
http://localhost:8081/api/v1
```

## Authentication

Currently, the API does not require authentication. This will be added in future phases.

## Endpoints

### Health Check

#### GET /health

Check the health status of the API server.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0"
}
```

### Phishlets

#### GET /phishlets

List all phishlets with optional filtering.

**Query Parameters:**
- `name` (string): Filter by phishlet name
- `enabled` (boolean): Filter by enabled status
- `limit` (integer): Maximum number of results
- `offset` (integer): Number of results to skip

**Response:**
```json
{
  "phishlets": [
    {
      "id": "linkedin",
      "name": "linkedin",
      "display_name": "LinkedIn",
      "author": "kgretzky",
      "version": "1.0.0",
      "description": "LinkedIn phishing template",
      "redirect_url": "https://linkedin.com",
      "proxy_hosts": [...],
      "is_enabled": true,
      "is_visible": true,
      "create_time": "2024-01-01T12:00:00Z",
      "update_time": "2024-01-01T12:00:00Z"
    }
  ],
  "count": 1
}
```

#### POST /phishlets

Create a new phishlet.

**Request Body:**
```json
{
  "name": "example",
  "display_name": "Example Site",
  "author": "user",
  "version": "1.0.0",
  "redirect_url": "https://example.com",
  "proxy_hosts": [...],
  "is_enabled": false
}
```

#### GET /phishlets/{name}

Get a specific phishlet by name.

#### PUT /phishlets/{name}

Update a specific phishlet.

#### DELETE /phishlets/{name}

Delete a specific phishlet.

#### GET /phishlets/{name}/stats

Get statistics for a specific phishlet.

**Response:**
```json
{
  "total_phishlets": 1,
  "enabled_phishlets": 1,
  "active_campaigns": 5
}
```

### Sessions

#### GET /sessions

List all sessions with optional filtering.

**Query Parameters:**
- `phishlet` (string): Filter by phishlet name
- `username` (string): Filter by username
- `start_time` (RFC3339): Filter by start time
- `end_time` (RFC3339): Filter by end time
- `limit` (integer): Maximum number of results
- `offset` (integer): Number of results to skip

**Response:**
```json
{
  "sessions": [
    {
      "id": "sess_123",
      "index": 1,
      "phishlet_name": "linkedin",
      "landing_url": "https://phish.example.com/login",
      "username": "victim@example.com",
      "password": "password123",
      "user_agent": "Mozilla/5.0...",
      "remote_addr": "192.168.1.100",
      "create_time": "2024-01-01T12:00:00Z",
      "update_time": "2024-01-01T12:05:00Z",
      "is_active": true
    }
  ],
  "count": 1
}
```

#### POST /sessions

Create a new session.

#### GET /sessions/{id}

Get a specific session by ID.

#### PUT /sessions/{id}

Update a specific session.

#### DELETE /sessions/{id}

Delete a specific session.

#### GET /sessions/stats

Get session statistics.

**Response:**
```json
{
  "total_sessions": 100,
  "active_sessions": 25,
  "captured_creds": 75,
  "unique_phishlets": 5
}
```

### Configuration

#### GET /config

Get the current configuration.

**Response:**
```json
{
  "general": {
    "domain": "example.com",
    "external_ipv4": "1.2.3.4",
    "bind_ipv4": "0.0.0.0",
    "https_port": 443,
    "dns_port": 53,
    "autocert": true
  },
  "proxy": {
    "type": "http",
    "address": "proxy.example.com",
    "port": 8080,
    "enabled": false
  },
  "update_time": "2024-01-01T12:00:00Z"
}
```

#### PUT /config

Update the configuration.

### Lures

#### GET /lures

List all lures.

#### POST /lures

Create a new lure.

#### GET /lures/{id}

Get a specific lure by ID.

#### PUT /lures/{id}

Update a specific lure.

#### DELETE /lures/{id}

Delete a specific lure.

### Certificates

#### GET /certificates

List all certificates.

**Response:**
```json
{
  "certificates": [
    {
      "domain": "example.com",
      "issuer": "Let's Encrypt",
      "not_before": "2024-01-01T00:00:00Z",
      "not_after": "2024-04-01T00:00:00Z",
      "is_valid": true,
      "is_wildcard": false
    }
  ],
  "count": 1
}
```

#### POST /certificates/generate

Generate a new certificate.

**Request Body:**
```json
{
  "domain": "example.com"
}
```

#### DELETE /certificates/{domain}

Delete a certificate for a specific domain.

## Error Responses

All endpoints return standard HTTP status codes. Error responses have the following format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common status codes:
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Usage Examples

### cURL Examples

```bash
# Check API health
curl http://localhost:8081/api/v1/health

# List all phishlets
curl http://localhost:8081/api/v1/phishlets

# Get specific phishlet
curl http://localhost:8081/api/v1/phishlets/linkedin

# List sessions with filtering
curl "http://localhost:8081/api/v1/sessions?phishlet=linkedin&limit=10"

# Get session statistics
curl http://localhost:8081/api/v1/sessions/stats

# Update configuration
curl -X PUT http://localhost:8081/api/v1/config \
  -H "Content-Type: application/json" \
  -d '{"general":{"domain":"newdomain.com"}}'
```

### Python Examples

```python
import requests
import json

# Base URL
base_url = "http://localhost:8081/api/v1"

# Check health
response = requests.get(f"{base_url}/health")
print(response.json())

# List phishlets
response = requests.get(f"{base_url}/phishlets")
phishlets = response.json()["phishlets"]

# Create a new session
session_data = {
    "id": "new_session_123",
    "phishlet_name": "linkedin",
    "landing_url": "https://phish.example.com/login",
    "user_agent": "Mozilla/5.0...",
    "remote_addr": "192.168.1.100"
}

response = requests.post(
    f"{base_url}/sessions",
    headers={"Content-Type": "application/json"},
    data=json.dumps(session_data)
)

if response.status_code == 201:
    print("Session created successfully")
else:
    print(f"Error: {response.json()['error']}")
```

## Integration with Legacy System

The API includes a legacy bridge that allows integration with the existing Evilginx2 core functionality. This ensures backward compatibility while providing modern API access.

## Future Enhancements

- Authentication and authorization
- WebSocket support for real-time updates
- Bulk operations
- Advanced filtering and search
- Rate limiting
- API versioning
- OpenAPI/Swagger documentation

## Development

To start the control service:

```bash
# Build the control service
go build -o control ./cmd/control

# Start with default settings
./control

# Start with custom port and database
./control -port 8082 -db ./custom.db

# Sync from legacy database
./control -sync -legacy-db ./legacy.db
```

The API server will start on the specified port (default: 8081) and provide all the endpoints documented above.
