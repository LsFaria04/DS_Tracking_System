# Google Cloud Storage Bucket for Blockchain Persistence

resource "google_storage_bucket" "blockchain_storage" {
  name          = "${var.blockchain_bucket_name}-${var.project_id}"
  location      = var.region
  force_destroy = false # Protect blockchain data from accidental deletion

  uniform_bucket_level_access = true

  versioning {
    enabled = true # Keep version history of blockchain snapshots
  }

  lifecycle_rule {
    condition {
      age = 90 # Archive old versions after 90 days
    }
    action {
      type          = "SetStorageClass"
      storage_class = "NEARLINE"
    }
  }

  labels = {
    environment = var.environment
    purpose     = "blockchain-storage"
  }
}

# Service account for Cloud Run to access blockchain storage
resource "google_service_account" "cloudrun_sa" {
  account_id   = "${var.service_name}-sa"
  display_name = "Service Account for ${var.service_name}"
  description  = "Service account for Cloud Run to access Cloud Storage and Cloud SQL"
}

# Grant Cloud Storage access to service account
resource "google_storage_bucket_iam_member" "blockchain_storage_access" {
  bucket = google_storage_bucket.blockchain_storage.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.cloudrun_sa.email}"
}

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
    service_account = google_service_account.cloudrun_sa.email

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

      # Blockchain Configuration
      env {
        name  = "BLOCKCHAIN_STORAGE_BUCKET"
        value = google_storage_bucket.blockchain_storage.name
      }
      env {
        name  = "BLOCKCHAIN_STORAGE_PATH"
        value = "blockchain/chain.json"
      }
      env {
        name  = "BLOCKCHAIN_DIFFICULTY"
        value = tostring(var.blockchain_difficulty)
      }
      env {
        name  = "BLOCKCHAIN_BACKUP_ENABLED"
        value = "true"
      }
      env {
        name  = "BLOCKCHAIN_BACKUP_INTERVAL"
        value = "300" # Backup every 5 minutes
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

output "blockchain_bucket_name" {
  description = "Name of the Cloud Storage bucket for blockchain data."
  value       = google_storage_bucket.blockchain_storage.name
}

output "blockchain_bucket_url" {
  description = "URL of the blockchain storage bucket."
  value       = google_storage_bucket.blockchain_storage.url
}

output "service_account_email" {
  description = "Email of the service account used by Cloud Run."
  value       = google_service_account.cloudrun_sa.email
}

output "environment_variables" {
  description = "Environment variables configured for the service."
  value = {
    DB_HOST                    = google_sql_database_instance.postgres.public_ip_address
    DB_PORT                    = "5432"
    DB_NAME                    = var.db_name
    BLOCKCHAIN_STORAGE_BUCKET  = google_storage_bucket.blockchain_storage.name
    BLOCKCHAIN_DIFFICULTY      = var.blockchain_difficulty
    ENVIRONMENT                = var.environment
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
    

      env {
        name = "API_URL"
        value = google_cloud_run_v2_service.default.uri
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