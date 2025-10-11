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
