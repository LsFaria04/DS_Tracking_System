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
  description = "Container image of backend to deploy to Cloud Run (including tag)."
  type        = string
  default     = "us-docker.pkg.dev/cloudrun/container/tracking-status"
}

variable "frontend_docker_image" {
  description = "Container image of frontend to deploy to Cloud Run (including tag)."
  type =  string
  default = "us-docker.pkg.dev/cloudrun/container/tracking-status-frontend"
  
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

variable "blockchain_bucket_name" {
  description = "Name for the Google Cloud Storage bucket to store blockchain data. Must be globally unique."
  type        = string
  default     = "tracking-blockchain-storage"
}

variable "environment" {
  description = "Environment name (development, staging, production)"
  type        = string
  default     = "production"
}

variable "blockchain_difficulty" {
  description = "Blockchain mining difficulty level (1-5). Higher = more secure but slower."
  type        = number
  default     = 2
  
  validation {
    condition     = var.blockchain_difficulty >= 1 && var.blockchain_difficulty <= 5
    error_message = "Blockchain difficulty must be between 1 and 5."
  }
}