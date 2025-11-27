# AI Coding Agent Instructions - Order Tracking System

## Architecture Overview

This is a distributed order tracking system with blockchain verification, deployed on GCP. It consists of:

- **Backend**: Go (1.25.3) REST API using Gin framework + GORM ORM + PostgreSQL
- **Frontend**: React 18.3.1 (standalone app using Vite)
- **Microfrontend**: `mips_order_tracking` (Module Federation with Rsbuild) exposing `OrderTrackingPage`
- **Blockchain**: Ethereum Sepolia testnet for immutable order status verification
- **Infrastructure**: GCP (Cloud Run + Cloud SQL) managed by Terraform
- **Event Bus**: Google Pub/Sub for inter-service communication

### Key Data Flow

1. Orders created → stored in PostgreSQL with GPS coordinates
2. Order status updates → dual-write to database + blockchain (as SHA-256 hash)
3. Verification → fetches blockchain hashes, compares with computed database hashes
4. Map visualization → uses Leaflet/React-Leaflet to render order trajectory with intermediate stops

## Critical Development Workflows

### Starting Development Environment

```bash
# Backend (includes PostgreSQL, PgAdmin, Pub/Sub emulator)
cd backend
docker compose -f compose.dev.yml up --build --watch

# Frontend (standalone)
cd frontend
docker compose -f compose.dev.yml up --build --watch
# Or: npm run dev

# Microfrontend (module federation)
cd mips_order_tracking
npm run dev
```

**Important**: Backend requires `.env` file in `backend/` with DB credentials + blockchain config (see `compose.dev.yml` env vars). Never commit `.env`!

### Database Management

- Migrations auto-run on container startup from `backend/db/migrations/`
- Seeds at `backend/db/seeds/dev_data.sql` must be manually run via PgAdmin (http://localhost:4321)
- Schema uses PostgreSQL enums (`order_state`) and triggers to prevent history updates/deletes
- Order price auto-calculated via trigger when order products change

### Testing Blockchain Integration

- Backend auto-checks Pub/Sub emulator on startup (`testPubSub()` in `main.go`)
- If `BLOCKCHAIN_RPC_URL` is empty, blockchain client returns `nil` (graceful degradation)
- Deploy contract: `GET /blockchain/deploy` (dev only, not for production)
- Check status: `GET /blockchain/status`
- Verify order: `GET /order/verify/:order_id` (compares DB hashes with blockchain)

### GCP Deployment

```bash
# 1. Build and push images
cd backend && docker build -f Dockerfile -t yourname/tracking-status:tag . && docker push ...
cd frontend && docker build -f Dockerfile -t yourname/tracking-status-frontend:tag . && docker push ...

# 2. Deploy (updates Cloud Run + Cloud SQL)
cd terraform
terraform apply -replace="google_cloud_run_v2_service.default"  # Backend
terraform apply -replace="google_cloud_run_v2_service.frontend" # Frontend

# 3. Make frontend public
gcloud run services add-iam-policy-binding tracking-status-frontend \
  --region=europe-west1 --member=allUsers --role=roles/run.invoker
```

**Or use**: `./deploy.ps1` (PowerShell script that automates above)

## Project-Specific Conventions

### Backend Handler Pattern

All handlers follow struct-based pattern with dependency injection:

```go
type OrderHandler struct {
    DB     *gorm.DB
    Client *blockchain.Client  // May be nil if blockchain not configured
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) { /* ... */ }
```

Handlers registered in `routes/router.go` via `RegisterRoutes(router, db, blockchainClient)`.

### Database Models (GORM)

- Field names use PascalCase with underscores (e.g., `Customer_ID`, `Tracking_Code`)
- Tags specify column types: `gorm:"type:decimal(10,8)"` for GPS coordinates
- Foreign keys use `foreignKey` tag, preload with `.Preload("Storage")`
- Example: `backend/src/models/order.go`

### Blockchain Integration

Order status updates stored as SHA-256 hash on Sepolia:

```go
data := fmt.Sprintf("%d|%s|%s|%s", order_id, status, timestamp, location)
hash := sha256.Sum256([]byte(data))
```

Contract ABI auto-generated: `blockchain/order_tracker_contract.go` (don't edit manually). Solidity source: `blockchain/contract.sol`.

### CORS Configuration

Backend allows:
- All `localhost:*` origins (dev)
- Any `*.run.app` domain (GCP Cloud Run)

See `main.go` `configRouter()` for `AllowOriginFunc` logic.

### Frontend Environment Variables

- Dev: `import.meta.env.VITE_API_URL` (defaults to `http://localhost:8080`)
- Production: Set via Docker build args or Terraform outputs

### Module Federation (Microfrontend)

`mips_order_tracking` exposes `./OrderTrackingPage` via `module-federation.config.ts`. Shared deps: React, React Router, Leaflet (all singleton, eager loading).

## Common Gotchas

1. **Database triggers**: `order_status_history` blocks UPDATE/DELETE via trigger. Only INSERT allowed.
2. **Blockchain nullable**: Always check `if h.Client != nil` before blockchain ops.
3. **Docker rebuilds**: When Compose/Dockerfile changes, run `down -v`, `build --no-cache`, `up --watch`.
4. **PgAdmin credentials**: `fakemail@gmail.com` / `safePassword` (from `compose.dev.yml`).
5. **GPS coordinates**: Use `DECIMAL(10,8)` for latitude, `DECIMAL(11,8)` for longitude.
6. **Terraform state**: Track `terraform.tfvars` (contains secrets) - don't commit!
7. **Air hot reload**: Dev backend uses Air (installed in `Dockerfile.dev`), watches `src/` directory.

## File Structure Highlights

- `backend/src/routes/router.go`: All API routes defined here
- `backend/db/migrations/001_init_schema.sql`: Schema with triggers
- `backend/src/blockchain/client.go`: Ethereum connection management
- `frontend/src/types.ts`: TypeScript interfaces for API responses
- `mips_order_tracking/module-federation.config.ts`: Microfrontend config
- `terraform/main.tf`: Complete GCP infrastructure (Cloud Run + Cloud SQL)

## Key Dependencies

- **Backend**: `gorm.io/gorm`, `github.com/gin-gonic/gin`, `github.com/ethereum/go-ethereum`, `cloud.google.com/go/pubsub`
- **Frontend**: React 18.3.1, React Router, Leaflet, TailwindCSS
- **Blockchain**: Sepolia testnet (Chain ID: 11155111) via Infura

## Integration Points

- **Jumpseller API**: Product data fetched from external API (env: `JUMPSELLER_BASE_URL`, `LOGIN_JUMPSELLER_API`, `TOKEN_JUMPSELLER_API`)
- **Pub/Sub**: Backend publishes to `orders` topic when `PUBSUB_EMULATOR_HOST` set (dev emulator at `http://localhost:8084`)
- **Checkout service**: Intended integration via Pub/Sub (not yet implemented per Sprint 1 notes)

## Testing/Debugging

- Backend logs: `gcloud run services logs read tracking-status --region=europe-west1` (production)
- Pub/Sub dashboard: http://localhost:8084 (dev emulator)
- PgAdmin: http://localhost:4321 (query execution, data inspection)
- Health check: `GET /ping` returns `{"message": "pong"}`
