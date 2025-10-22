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

      env {
        name  = "DB_HOST"
        value = google_sql_database_instance.postgres.public_ip_address
      }
      env {
        name  = "DB_USER"
        value = var.db_user
      }
      env {
        name  = "DB_PASSWORD"
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