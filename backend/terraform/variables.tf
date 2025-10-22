variable "project_id" {
  description = "The Google Cloud project ID."
  type        = string
}

variable "region" {
  description = "The region where resources will be deployed."
  type        = string
  default     = "europe-west1"
}

variable "service_name" {
  description = "The name for the Cloud Run service."
  type        = string
  default     = "tracking-status"
}

variable "docker_image" {
  description = "Container image to deploy to Cloud Run (including tag)."
  type        = string
  default     = "us-docker.pkg.dev/cloudrun/container/tracking-status"
}

variable "db_instance_name" {
  description = "The name for the Cloud SQL instance."
  type        = string
  default     = "tracking-db"
}

variable "db_name" {
  description = "The name of the PostgreSQL database to create."
  type        = string
  default     = "tracking_db"
}

variable "db_user" {
  description = "The username for the PostgreSQL database."
  type        = string
  default     = "tracking_user"
}

variable "db_password" {
  description = "The password for the PostgreSQL database user. IMPORTANT: Use a strong password and consider using environment variables or secret managers."
  type        = string
  sensitive   = true
}