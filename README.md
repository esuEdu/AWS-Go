# AWS Go API with RDS (Terraform on EC2)

This repository contains a minimal Go REST API (CRUD for Items) backed by PostgreSQL on Amazon RDS, plus Terraform to provision the required AWS resources: EC2, RDS, and IAM roles.

## API

-   Endpoints

    -   GET /items
    -   POST /items
    -   GET /items/{id}
    -   PUT /items/{id}
    -   DELETE /items/{id}

-   Environment variables
    -   ADDR (default ":8080")
    -   AWS_REGION
    -   DB_HOST, DB_PORT (default 5432), DB_USER, DB_NAME
    -   DB_PASSWORD or DB_PASSWORD_SSM_PARAM (name of a SecureString in SSM Parameter Store)

### Run locally

Prerequisites: Go 1.21+ and a reachable PostgreSQL database.

```
export ADDR=":8080"
export AWS_REGION=us-east-1
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=appuser
export DB_NAME=appdb
export DB_PASSWORD=yourpass

go run ./cmd/api
```

## Terraform

The Terraform in `infra/` provisions:

-   RDS (PostgreSQL 16) in the default VPC
-   EC2 (Amazon Linux 2023) with an instance profile/role to read the DB password from SSM
-   Security groups to allow the EC2 instance to talk to RDS, and expose port 8080 for the API
-   An SSM Parameter (SecureString) `/app/db_password` holding the generated DB password

The EC2 instance bootstraps itself to:

-   Install Go and git
-   Clone this repo, build the API, and run it as a systemd service on port 8080

### Deploy

```
cd infra
terraform init
terraform apply -auto-approve
```

Outputs include `api_url` and `rds_endpoint`.

To customize:

-   `var.region` (default `us-east-1`)
-   `var.db_username` (default `appuser`)
-   `var.db_name` (default `appdb`)
-   `var.allow_cidr_http` (default `0.0.0.0/0`)

### Clean up

```
cd infra
terraform destroy
```

## Project structure

-   `cmd/api/main.go` — server entry point
-   `internal/config` — config and SSM password fetch
-   `internal/repository/postgres` — DB pool and migrations
-   `internal/handlers` — HTTP handlers
-   `infra/` — Terraform (EC2, IAM, RDS)
