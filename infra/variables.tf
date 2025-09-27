variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "db_username" {
  description = "Master username for RDS"
  type        = string
  default     = "appuser"
}

variable "db_name" {
  description = "Initial database name"
  type        = string
  default     = "appdb"
}

variable "allow_cidr_http" {
  description = "CIDR block allowed to reach the API on port 8080"
  type        = string
  default     = "0.0.0.0/0"
}
