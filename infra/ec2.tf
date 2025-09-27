resource "aws_security_group" "ec2" {
  name        = "app-ec2-sg"
  description = "EC2 SG for app"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "HTTP"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [var.allow_cidr_http]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Latest Amazon Linux 2023 x86_64

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["137112412989"] # Amazon

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

locals {
  user_data = <<-EOF
              #!/bin/bash
              set -euo pipefail
              dnf update -y
              dnf install -y golang git

              useradd -m -s /bin/bash app || true
              install -d -o app -g app /opt/app

              # Fetch DB password from SSM
              DB_PASSWORD=$(aws ssm get-parameter --name "/app/db_password" --with-decryption --query Parameter.Value --output text --region ${var.region})

              cat >/etc/goapp.env <<EOT
              ADDR=":8080"
              AWS_REGION="${var.region}"
              DB_HOST="${aws_db_instance.this.address}"
              DB_PORT="5432"
              DB_USER="${var.db_username}"
              DB_NAME="${var.db_name}"
              DB_PASSWORD="$DB_PASSWORD"
              EOT

              su - app -c "git clone https://github.com/esuEdu/AWS-Go /home/app/app || (cd /home/app/app && git pull)"
              cd /home/app/app
              su - app -c "cd /home/app/app && /usr/bin/go build -o /opt/app/app ./cmd/api"

              cat >/etc/systemd/system/goapp.service <<EOT
              [Unit]
              Description=Go API
              After=network-online.target

              [Service]
              EnvironmentFile=/etc/goapp.env
              ExecStart=/opt/app/app
              User=app
              WorkingDirectory=/opt/app
              Restart=always
              RestartSec=3

              [Install]
              WantedBy=multi-user.target
              EOT

              systemctl daemon-reload
              systemctl enable --now goapp
              EOF
}

resource "aws_instance" "app" {
  ami                         = data.aws_ami.al2023.id
  instance_type               = "t3.micro"
  subnet_id                   = data.aws_subnets.default.ids[0]
  vpc_security_group_ids      = [aws_security_group.ec2.id]
  iam_instance_profile        = aws_iam_instance_profile.ec2.name
  user_data                   = local.user_data

  metadata_options {
    http_tokens = "required"
  }

  tags = {
    Name = "go-api"
  }
}
