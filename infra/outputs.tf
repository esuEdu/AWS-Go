output "api_url" {
  description = "Public URL for the API"
  value       = "http://${aws_instance.app.public_ip}:8080"
}

output "rds_endpoint" {
  value       = aws_db_instance.this.address
  description = "RDS endpoint hostname"
}

output "db_password_param" {
  value       = aws_ssm_parameter.db_password.name
  description = "SSM parameter name that stores DB password"
}
