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

Certificates (stub)
- GET /certificates
- POST /certificates
- DELETE /certificates/:domain

Example
curl -H "Authorization: Bearer $API_ADMIN_TOKEN" http://localhost:8081/api/v1/phishlets
