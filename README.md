# Order Tracking System

Backend developed with Go, frontend with React (v18.0.0), deployed on Google Cloud Platform.

The main development milestones and changes are reported in the [changelog](docs/CHANGELOG.MD) file.

## API Documentation

For detailed information, refer to the backend API documentation available [here](https://documenter.getpostman.com/view/44004544/2sB3dPTVyi).

## Quick Start

To start developing use the following command to create the docker containers:

```shell
# Use the dev composer file to start postgres, pgadmin and the go containers
cd backend
docker compose -f compose.dev.yml up  --build --watch # The build flag should only be used when changes are made when the container is not running (ex : pulling code from git)
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

To run the Compose and Dockerfile files, make sure to create a .env file containing the required environment variables. Place this file in the same directory as your Compose configuration to ensure proper loading.

**DO NOT COMMIT THE .ENV FILE!**

## Tracking system frontend

## Quick Start

```bash
# From project root
cd mips_order_tracking

# Copy environment template and install deps
cp .env.example .env #Some keys may need to be included manually
npm install

# Start dev container with hot reload (preferred)
docker compose -f compose.dev.yml up --build --watch  # The build flag should only be used when changes are made when the container is not running (ex : pulling code from git)
# Or run without Docker:
npm run dev
```

### Accessing Services

- **Frontend**: http://localhost:5174
- **Backend**: http://localhost:8080 (must be running for API calls)

To stop the frontend containers:

```bash
cd mips_order_tracking
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

# Navigate to terraform directory
cd terraform

# Create configuration file
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your:
#   - project_id
#   - db_password
#   - blockchain_rpc_url (Infura Sepolia endpoint)
#   - blockchain_private_key (MetaMask private key - do not share this)
#   - blockchain_contract_address (leave empty until contract is deployed)

# Initialize Terraform (first time only)
terraform init
```

### Get MetaMask Private Key

**SECURITY WARNING**: Never commit your private key

1. Open MetaMask browser extension
2. Click the 3 dots menu on your account
3. Select "Account details"
4. Click "Show private key"
5. Enter your MetaMask password
6. Copy the private key
7. Paste it in `terraform/terraform.tfvars` under `blockchain_private_key`

### Deploy to GCP

```shell
# Build and push Docker image backend
cd backend
docker build -f Dockerfile -t <yourname>/tracking-status:<tag> .
docker push <yourname>/tracking-status:<tag>

# Build and push Docker image frontend 
cd ../mips_order_tracking
docker build -f Dockerfile -t <yourname>/tracking-status-frontend:<tag> .
docker push <yourname>/tracking-status-frontend:<tag>

# Deploy Cloud SQL + Cloud Run + Blockchain Config
cd ../terraform
terraform apply -replace="google_cloud_run_v2_service.default"  # To ensure backend updates

terraform apply -replace="google_cloud_run_v2_service.frontend"  # To ensure frontend updates (if run this command, you need to bind the iam policy to be able to enter the app)

# IAM policy for frontend public access
gcloud run services add-iam-policy-binding tracking-status-frontend --region=europe-west1 --member=allUsers --role=roles/run.invoker
```

**Alternative: Use deployment script**

PowerShell (Windows):
```powershell
.\deploy.ps1
```

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
