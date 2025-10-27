# Order Tracking System

Backend developed with Go, frontend with Next.js, deployed on Google Cloud Platform.

## Quick Start

To start developing use the following command to create the docker containers:

```shell
# Use the dev composer file to start postgres, pgadmin and the go containers
cd backend
docker compose -f compose.dev.yml up --watch
```

The database schema will be automatically initialized on first run. However, to populate with data the development database you need to access pgadmin and run the data seed provided in db/seeds.

### Accessing Services

- **Backend API**: http://localhost:8080
- **PgAdmin**: http://localhost:4321

To stop the containers and remove them use:

```shell
docker compose -f compose.dev.yml  down
```

Whenever changes are made to a Dockerfile or Compose file, use the following commands to rebuild and ensure those updates  are applied to your containers:

```shell
docker compose -f compose.dev.yml down -v
docker compose -f compose.dev.yml build --no-cache
docker compose -f compose.dev.yml up --watch
```

To run the production Compose and Dockerfile (i.e., the ones without the .dev suffix), make sure to create a .env file containing the required environment variables. Place this file in the same directory as your Compose configuration to ensure proper loading.

**DO NOT COMMIT THE .ENV FILE!**

## Tracking system frontend

## Quick Start

```bash
# From project root
cd frontend

# Copy environment template and install deps
cp .env.example .env
npm install

# Start dev container with hot reload (preferred)
docker compose -f compose.dev.yml up --watch
# Or run without Docker:
npm run dev
```

### Accessing Services

- **Frontend**: http://localhost:3000
- **Backend**: http://localhost:8080 (must be running for API calls)

To stop the frontend containers:

```bash
cd frontend
docker compose -f compose.dev.yml down
```

## Terraform Deployment to GCP

### Prerequisites
1. Install [Terraform](https://www.terraform.io/downloads)
2. Install and authenticate [gcloud CLI](https://cloud.google.com/sdk/docs/install)

### Configuration

Terraform manages environment variables automatically. For local development, edit `backend/.env`:

```bash
DB_HOST=postgres
DB_PORT=5432
DB_USER=tracking_user
DB_PASS=change_me
DB_NAME=tracking_db
```

### Setup

```shell
# Authenticate with GCP
gcloud auth application-default login
gcloud config set project madeinportugal

# Enable required APIs
gcloud services enable cloudrun.googleapis.com
gcloud services enable sqladmin.googleapis.com
gcloud services enable storage.googleapis.com

# Navigate to terraform directory
cd backend/terraform

# Create configuration file
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your project_id, db_password, etc.

# Initialize Terraform (first time only)
terraform init
```

### Deploy to GCP

```shell
# Build and push Docker image
cd backend
docker build -f Dockerfile -t returnedft/tracking-status:test .
docker push returnedft/tracking-status:test

# Deploy Cloud SQL + Cloud Run
cd terraform
terraform apply
```

Takes 10-15 minutes. Type `yes` when prompted.

Get your service URL:
```shell
terraform output service_url
```

### Stop Resources (Stop Charges)

```shell
cd backend/terraform
terraform destroy
```

Type `yes` when prompted.

**Important**: `terraform destroy` deletes:
- Cloud SQL instance (including all data)
- Cloud Run service
- All associated resources

This action is irreversible but stops all charges immediately.

### Daily Workflow

```shell
# Morning: Deploy for testing
cd backend/terraform
terraform apply

# Evening: Destroy to save credits
terraform destroy
```

### View Logs

```shell
gcloud run services logs read tracking-status --region=europe-west1
```

### Troubleshooting

**Check Cloud SQL status:**
```shell
gcloud sql instances describe tracking-db
```

**Check Cloud Run status:**
```shell
gcloud run services describe tracking-status --region=europe-west1
```

**If APIs aren't enabled:**
```shell
gcloud services enable cloudrun.googleapis.com sqladmin.googleapis.com storage.googleapis.com
```