# Nimble Fullstack Challenge (Pet Store)

This is a small, local pet store app built for the Nimble exercise. I kept it intentionally simple but production‑ready in structure: GraphQL API in Go, React (TypeScript) UI for customers, Postgres via Docker. Everything runs locally and doesn’t depend on external services.

## Quick start (Docker)

Requirements: Docker Desktop (or Docker Engine + Compose).

```
docker compose up --build
```

Open the customer UI:

```
http://localhost:3000/store/demo
```

The API is HTTPS with a self‑signed cert. The first time your browser hits it, you’ll need to accept the warning:

```
https://localhost:8443/graphql
```

## Local URLs

- UI: `http://localhost:3000/store/demo`
- API: `https://localhost:8443/graphql`
- Postgres: `localhost:5432` (credentials in `.env`)

## Auth (Basic)

Both customer and merchant endpoints are protected with Basic Auth.

- Merchant: `merchant_demo` / `merchant_demo_pw`
- Customer: `customer_demo` / `customer_demo_pw`

Credentials can be changed in `.env`. Demo store + users seed on API startup.

## What’s in the repo

- `backend/`: Go GraphQL server, schema, seed logic
- `frontend/`: React + TypeScript UI
- `infra/`: Postgres schema and local infra bits
- `docker-compose.yml`: local orchestration

## A few GraphQL examples

Create a pet (merchant):

```
curl -k -u merchant_demo:merchant_demo_pw \
  -H "Content-Type: application/json" \
  https://localhost:8443/graphql \
  -d '{
    "query":"mutation($input: CreatePetInput!){ createPet(input:$input){ id name species createdAt } }",
    "variables":{
      "input":{
        "name":"Miso",
        "species":"CAT",
        "ageYears":2,
        "pictureUrl":"https://example.com/miso.jpg",
        "description":"Playful kitten",
        "breederName":"Jane Doe",
        "breederEmail":"jane@example.com"
      }
    }
  }'
```

Purchase pets (customer):

```
curl -k -u customer_demo:customer_demo_pw \
  -H "Content-Type: application/json" \
  https://localhost:8443/graphql \
  -d '{
    "query":"mutation($input: PurchasePetsInput!){ purchasePets(input:$input){ purchasedIds errors{ petName message } } }",
    "variables":{
      "input":{
        "storeSlug":"demo",
        "petIds":["pet-id-1","pet-id-2"]
      }
    }
  }'
```

## UI features

- Store page shows available pets only
- Cart + checkout (all at once)
- Error message if pets were already purchased
- “Add item” tab for quick merchant testing
- “History” tab showing purchased pets

## Security notes (short version)

- Passwords are hashed with Argon2id
- Breeder emails are encrypted at rest (AES‑GCM)
- Purchases are transactional with row locks (`SELECT … FOR UPDATE`)
- Basic rate limiting and safe headers on the API

## Optional dev (no Docker)

Backend:

```
cd backend
go run ./cmd/api
```

Frontend:

```
cd frontend
npm install
npm run dev
```

## Gotchas / tips

- If the UI says “Failed to fetch”, it’s almost always the self‑signed cert. Visit the API URL once and accept the warning.
- If you want a clean database, remove `infra/db-data` and restart compose.
