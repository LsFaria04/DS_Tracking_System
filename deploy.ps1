# Backend
Write-Host "Building backend..." -ForegroundColor Yellow
Set-Location backend
docker build -t returnedft/tracking-status:test .
docker push returnedft/tracking-status:test

# Frontend  
Write-Host "Building frontend..." -ForegroundColor Yellow
Set-Location ..\frontend
docker build -t returnedft/tracking-status-frontend:testfrontend .
docker push returnedft/tracking-status-frontend:testfrontend

# Terraform
Write-Host "Running terraform..." -ForegroundColor Yellow
Set-Location ..\terraform
terraform apply

# Fix IAM
Write-Host "Fixing IAM..." -ForegroundColor Yellow
gcloud run services add-iam-policy-binding tracking-status-frontend `
  --region=europe-west1 `
  --member=allUsers `
  --role=roles/run.invoker

Write-Host "Done!" -ForegroundColor Green
terraform output service_url
terraform output frontend_url
