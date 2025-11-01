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

variable "blockchain_rpc_url" {
  description = "Ethereum RPC URL for connecting to Sepolia testnet (e.g., Infura endpoint)"
  type        = string
  sensitive   = true
}

variable "blockchain_private_key" {
  description = "Private key for Ethereum wallet to sign transactions. NEVER commit this to version control!"
  type        = string
  sensitive   = true
}

variable "blockchain_contract_address" {
  description = "Deployed smart contract address on Sepolia testnet (leave empty if not yet deployed)"
  type        = string
  default     = "0x472b477d30c45cfbd89e76d9f9700ad1f90cc370"
}