# Generate unique tag based on timestamp
$TAG = Get-Date -Format "yyyyMMdd-HHmmss"

# Backend
Write-Host "Building backend..." -ForegroundColor Yellow
Set-Location backend
docker build -t returnedft/tracking-status:$TAG .
docker push returnedft/tracking-status:$TAG

# Frontend  
Write-Host "Building frontend..." -ForegroundColor Yellow
Set-Location ..\mips_order_tracking
docker build --build-arg PUBLIC_API_URL=https://tracking-status-138871440259.europe-west1.run.app --build-arg PUBLIC_TOMTOM_API_KEY=$env:PUBLIC_TOMTOM_API_KEY -t returnedft/tracking-status-frontend:$TAG .
docker push returnedft/tracking-status-frontend:$TAG


# Terraform
Write-Host "Running terraform with new images..." -ForegroundColor Yellow
Set-Location ..\terraform

# Import existing resources if they exist
Write-Host "Importing existing resources..." -ForegroundColor Cyan
terraform import -input=false google_sql_database_instance.postgres madeinportugal/tracking-db 2>$null
terraform import -input=false google_sql_database.database madeinportugal/tracking-db/tracking_db 2>$null
terraform import -input=false google_sql_user.user madeinportugal/tracking-db/tracking_user 2>$null
terraform import -input=false google_cloud_run_v2_service.default projects/madeinportugal/locations/europe-west1/services/tracking-status 2>$null
terraform import -input=false google_cloud_run_v2_service.frontend projects/madeinportugal/locations/europe-west1/services/tracking-status-frontend 2>$null

Write-Host "Applying Terraform changes..." -ForegroundColor Yellow
terraform apply -auto-approve `
  -var="docker_image=returnedft/tracking-status:$TAG" `
  -var="frontend_docker_image=returnedft/tracking-status-frontend:$TAG" `
  -var="blockchain_contract_address=0xCB0B5282057FCf183dE89CF3115a01a02e82eB61"


# Fix IAM
Write-Host "Fixing IAM..." -ForegroundColor Yellow
gcloud run services add-iam-policy-binding tracking-status-frontend `
  --region=europe-west1 `
  --member=allUsers `
  --role=roles/run.invoker

Write-Host "Done!" -ForegroundColor Green
terraform output service_url
terraform output frontend_url
