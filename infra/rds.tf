resource "random_password" "db" {
  length  = 20
  special = true
}

resource "aws_db_subnet_group" "this" {
  name       = "app-default-subnets"
  subnet_ids = data.aws_subnets.default.ids
}

resource "aws_security_group" "rds" {
  name        = "app-rds-sg"
  description = "Allow Postgres from EC2"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "PostgreSQL from EC2 SG"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ec2.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "this" {
  identifier                  = "app-postgres"
  engine                      = "postgres"
  engine_version              = "16.4"
  instance_class              = "db.t3.micro"
  allocated_storage           = 20
  db_subnet_group_name        = aws_db_subnet_group.this.name
  vpc_security_group_ids      = [aws_security_group.rds.id]
  username                    = var.db_username
  password                    = random_password.db.result
  db_name                     = var.db_name
  skip_final_snapshot         = true
  publicly_accessible         = false
  deletion_protection         = false
  performance_insights_enabled = false
}

resource "aws_ssm_parameter" "db_password" {
  name  = "/app/db_password"
  type  = "SecureString"
  value = random_password.db.result
}
