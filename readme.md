# cloud-container-apis

Graphql APIs to spawn up a container & delete a container

# Getting started

### Required services

- Postgres DB
- Redis
- SMTP Server
- Container Orchestrator

Sample `.env` file that can be used

```.env
SMTP_HOST=smtp.ethereal.email
SMTP_PORT=587
SMTP_USERNAME=nils.powlowski30@ethereal.email
SMTP_PASSWORD=864fhYNmEDV5wCVXzQ
SMTP_SENDER_EMAIL=no-reply@cloudcontainer.com
SMTP_SENDER_NAME=Lakhan Samani
CONTAINER_ORCHESTRATOR_SERVICE_URL=0.0.0.0:5600
REDIS_URL=redis://localhost:6379
DATABASE_URL=postgres://postgres:password@0.0.0.0:5432/postgres
```

### Running

- `make run`

### Building

- `make`
