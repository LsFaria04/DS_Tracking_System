# Cloud SQL PostgreSQL Instance
resource "google_sql_database_instance" "postgres" {
  name             = var.db_instance_name
  database_version = "POSTGRES_17"
  region           = var.region

  settings {
    tier = "db-f1-micro" # Smallest instance (~$10/month)

    backup_configuration {
      enabled = false # Disable backups to save costs during testing
    }

    ip_configuration {
      ipv4_enabled = true # Enable public IP for Cloud Run connectivity
      authorized_networks {
        name  = "allow-all"
        value = "0.0.0.0/0" # Allow all IPs (Cloud Run uses dynamic IPs)
      }
    }
  }

  deletion_protection = false # Allow terraform destroy without manual intervention
}

# Database
resource "google_sql_database" "database" {
  name     = var.db_name
  instance = google_sql_database_instance.postgres.name
}

# Database User
resource "google_sql_user" "user" {
  name     = var.db_user
  instance = google_sql_database_instance.postgres.name
  password = var.db_password
}

# Cloud Run Service
resource "google_cloud_run_v2_service" "default" {
  name     = var.service_name
  location = var.region

  template {
    containers {
      image = var.docker_image

      # Database Configuration
      env {
        name  = "DB_HOST"
        value = google_sql_database_instance.postgres.public_ip_address
      }
      env {
        name  = "DB_USER"
        value = var.db_user
      }
      env {
        name  = "DB_PASS"
        value = var.db_password
      }
      env {
        name  = "DB_NAME"
        value = var.db_name
      }
      env {
        name  = "DB_PORT"
        value = "5432"
      }

      # Blockchain Configuration (Sepolia Testnet)
      env {
        name  = "BLOCKCHAIN_RPC_URL"
        value = var.blockchain_rpc_url
      }
      env {
        name  = "BLOCKCHAIN_PRIVATE_KEY"
        value = var.blockchain_private_key
      }
      env {
        name  = "BLOCKCHAIN_CONTRACT_ADDRESS"
        value = var.blockchain_contract_address
      }
      env {
        name  = "BLOCKCHAIN_NETWORK"
        value = "sepolia"
      }

      # Application Configuration
      env {
        name  = "ENVIRONMENT"
        value = var.environment
      }
      env {
        name  = "GCP_PROJECT_ID"
        value = var.project_id
      }

      # Jumpseller API Credentials
      env {
        name  = "JUMPSELLER_BASE_URL"
        value = var.jumpseller_base_url
      }
      env {
        name  = "LOGIN_JUMPSELLER_API"
        value = var.login_jumpseller_api
      }
      env {
        name  = "TOKEN_JUMPSELLER_API"
        value = var.token_jumpseller_api
      }

      # Resource limits
      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }

    scaling {
      min_instance_count = 0
      max_instance_count = 10
    }
  }
}

resource "google_cloud_run_v2_service_iam_binding" "public_access" {
  project  = google_cloud_run_v2_service.default.project
  location = google_cloud_run_v2_service.default.location
  name     = google_cloud_run_v2_service.default.name

  role    = "roles/run.invoker"
  members = ["allUsers"]
}

output "service_url" {
  description = "The URL of the deployed Cloud Run service."
  value       = google_cloud_run_v2_service.default.uri
}

output "db_connection_name" {
  description = "Cloud SQL instance connection name."
  value       = google_sql_database_instance.postgres.connection_name
}

output "db_public_ip" {
  description = "Database public IP address used by Cloud Run."
  value       = google_sql_database_instance.postgres.public_ip_address
}

output "blockchain_config" {
  description = "Blockchain configuration for Sepolia testnet."
  value = {
    BLOCKCHAIN_NETWORK = "sepolia"
    BLOCKCHAIN_RPC_URL = "configured"
  }
  sensitive = false
}


# Frontend

resource "google_cloud_run_v2_service" "frontend" {
  name = "tracking-status-frontend"
  location = var.region

  template {
    containers {
      image = var.frontend_docker_image
    
      ports {
        container_port = 80
      }
    }


  }
  
}

# Frontend public access
resource "google_cloud_run_v2_service_iam_binding" "frontend_public_access" {
  project  = google_cloud_run_v2_service.frontend.project
  location = google_cloud_run_v2_service.frontend.location
  name     = google_cloud_run_v2_service.frontend.name

  role    = "roles/run.invoker"
  members = ["allUsers"]
}

output "frontend_url" {
  description = "The URL of the deployed frontend Cloud Run service."
  value       = google_cloud_run_v2_service.frontend.uri
}