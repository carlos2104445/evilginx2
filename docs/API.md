# Evilginx2 REST API

Base URL
- Default: http://localhost:8081/api/v1
- Auth: Bearer token via Authorization header (API_ADMIN_TOKEN)

Health
- GET /health

Config
- GET /config
- PUT /config

Phishlets
- GET /phishlets
- POST /phishlets
- GET /phishlets/:name
- PUT /phishlets/:name
- DELETE /phishlets/:name
- GET /phishlets/:name/stats

Advanced Phishlets
- GET /phishlets/:name/versions
  - Response: { "versions": [PhishletVersion], "count": n }
- POST /phishlets/:name/versions
  - Body: { "version": "v1", "description": "desc" }
  - Response: { "message": "version created" }
- GET /phishlets/:name/versions/:version
  - Response: models.Phishlet for the specified version

- POST /phishlets/:name/evaluate
  - Body: {
      "email": "user@corp.com",
      "user_agent": "UA",
      "ip_address": "1.2.3.4",
      "hostname": "sub.example.com",
      "path": "/login",
      "custom_params": { "k": "v" }
    }
  - Response: { "actions": [{ "type": "...", "value": "..." }...] }

- GET /phishlets/:name/flows
  - Response: { "flows": [...] } mapped from phishlet.MultiPageFlows

- POST /phishlets/:name/flows/:flow/step
  - Body: { "session_id": "sid", "step_data": { "k": "v" } }
  - Response: 201 with { "message": "flow session created" } on first call; 200 on update with { "message": "flow step updated" }

Sessions
- GET /sessions
- POST /sessions
- GET /sessions/:id
- PUT /sessions/:id
- DELETE /sessions/:id
- GET /sessions/stats

Lures
- GET /lures
- POST /lures
- GET /lures/:id
- PUT /lures/:id
- DELETE /lures/:id

Certificates
- GET /certificates
  - Response: { "certificates": [ { "domain": "...", ... } ], "count": n }
- POST /certificates
  - Body: {
      "domain": "test.example",
      "issuer": "Test CA",
      "not_before": "2025-01-01T00:00:00Z",
      "not_after": "2026-01-01T00:00:00Z",
      "is_wildcard": false
    }
  - Response: certificate object
- DELETE /certificates/:domain
  - Response: { "message": "certificate deleted successfully" }

# Evilginx2 REST API

Base URL: /api/v1

Auth
- Header: Authorization: Bearer <API_ADMIN_TOKEN>
- Unauthenticated: GET /health

CORS
- Allowed origin via FRONTEND_ORIGIN (default http://localhost:5173)

Health
- GET /health

Config
- GET /config
- PUT /config

Phishlets
- GET /phishlets
- POST /phishlets
- GET /phishlets/:name
- PUT /phishlets/:name
- DELETE /phishlets/:name
- GET /phishlets/:name/stats

Sessions
- GET /sessions
- POST /sessions
- GET /sessions/:id
- PUT /sessions/:id
- DELETE /sessions/:id
- GET /sessions/stats

Lures
- GET /lures
- POST /lures
- GET /lures/:id
- PUT /lures/:id
- DELETE /lures/:id

Phishlet Advanced
- GET /phishlets/:name/versions
- POST /phishlets/:name/versions
- GET /phishlets/:name/versions/:version
- POST /phishlets/:name/evaluate
- GET /phishlets/:name/flows
- POST /phishlets/:name/flows/:flow/step

Certificates
- GET /certificates
- POST /certificates
- DELETE /certificates/:domain

Examples

List phishlets:
curl -H "Authorization: Bearer $API_ADMIN_TOKEN" http://localhost:8081/api/v1/phishlets

Certificates:
- List
curl -H "Authorization: Bearer $API_ADMIN_TOKEN" http://localhost:8081/api/v1/certificates

- Create
curl -X POST -H "Authorization: Bearer $API_ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"domain":"test.example","issuer":"Test CA","not_before":"2025-01-01T00:00:00Z","not_after":"2026-01-01T00:00:00Z","is_wildcard":false}' \
  http://localhost:8081/api/v1/certificates

- Delete
curl -X DELETE -H "Authorization: Bearer $API_ADMIN_TOKEN" \
  http://localhost:8081/api/v1/certificates/test.example
