terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.54.1"
    }
  }
}
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support = true

  tags = {
    App = "movie-reviews"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "main-igw"
    App = var.AppName
  }
}

resource "aws_route_table" "main" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "main-route-table"
    App = var.AppName
  }
}

resource "aws_subnet" "subnet-1" {
  vpc_id = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
  availability_zone = "eu-central-1a"

  tags = {
    Name = "subnet-1"
    App = var.AppName
  }
}

resource "aws_subnet" "subnet-2" {
  vpc_id = aws_vpc.main.id
  cidr_block = "10.0.2.0/24"
  availability_zone = "eu-central-1b"

  tags = {
    Name = "subnet-2"
    App = var.AppName
  }
}

resource "aws_route_table_association" "a-1" {
  subnet_id = aws_subnet.subnet-1.id
  route_table_id = aws_route_table.main.id
}

resource "aws_route_table_association" "a-2" {
  subnet_id = aws_subnet.subnet-2.id
  route_table_id = aws_route_table.main.id
}

resource "aws_db_subnet_group" "main" {
  name = "main-db-subnet-group"
  subnet_ids = [aws_subnet.subnet-1.id, aws_subnet.subnet-2.id]

  tags = {
    App = var.AppName
  }
}

resource "aws_security_group" "allow-postgres" {
  name = "allow-postgres"
  description = "Allow inbound traffic for postgres and my personal computer to RDS"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port = 5432
    to_port = 5432
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow-postgres"
    App = var.AppName
  }
}

resource "aws_security_group" "allow-web" {
  name = "allow-web"
  description = "Allow inbound traffic for port 80"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port = 80
    to_port = 90
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow-web"
    App = var.AppName
  }
}

resource "aws_security_group" "allow-ssh" {
  name = "allow-ssh"
  description = "Allow ssh inbound traffic"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["1.2.3.0/24"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "allow-ssh"
    App = var.AppName
  }
}

resource "aws_iam_role" "ec2-role" {
  name = "movie-reviews-ec2-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "ec2-role-policy" {
  name   = "movie-reviews-ec2-role-policy"
  role   = aws_iam_role.ec2-role.id
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParameter"
      ],
      "Resource": "arn:aws:ssm:*:*:parameter/movie-reviews/*"
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "ec2-instance-profile" {
  name = "movie-reviews-ec2-instance-profile"
  role = aws_iam_role.ec2-role.name
}

resource "aws_instance" "host_instance" {
  ami = var.ami_id
  instance_type = var.instance_type
  key_name = "movie-reviews"
  iam_instance_profile = aws_iam_instance_profile.ec2-instance-profile.name
  vpc_security_group_ids = [aws_security_group.allow-web.id, aws_security_group.allow-ssh.id]
  user_data_replace_on_change = true
  subnet_id = aws_subnet.subnet-1.id
  associate_public_ip_address = true

  user_data = <<-EOF
    #!/bin/bash
    sudo apt-get update
    sudo apt-get install -y docker.io awscli jq
    sudo systemctl start docker
    sudo systemctl enable docker

    REGION=eu-central-1
    JWT_SECRET=$(aws ssm get-parameter --region $REGION --name "/movie-reviews/jwt-secret" --with-decryption --output json | jq -r Parameter.Value)
    ADMIN_NAME=$(aws ssm get-parameter --region $REGION --name "/movie-reviews/admin/name" --with-decryption --output json | jq -r Parameter.Value)
    ADMIN_EMAIL=$(aws ssm get-parameter --region $REGION --name "/movie-reviews/admin/email" --with-decryption --output json | jq -r Parameter.Value)
    ADMIN_PASSWORD=$(aws ssm get-parameter --region $REGION --name "/movie-reviews/admin/password" --with-decryption --output json | jq -r Parameter.Value)
    DB_URL=$(aws ssm get-parameter --region $REGION --name "/movie-reviews/db-url" --with-decryption --output json | jq -r Parameter.Value)

    sudo docker run -d \
      --name movie-reviews \
      -p 80:8000 \
      -e JWT_SECRET=$JWT_SECRET \
      -e ADMIN_NAME=$ADMIN_NAME \
      -e ADMIN_EMAIL=$ADMIN_EMAIL \
      -e ADMIN_PASSWORD=$ADMIN_PASSWORD \
      -e DB_URL=$DB_URL \
      davidmovas/movie-reviews:latest

    sudo docker run -d \
      --name watchtower \
      -v /var/run/docker.sock:/var/run/docker.sock \
      containrrr/watchtower
      davidmovas/movie-reviews:latest
      --schedule "0/30 * * * * *"
    EOF

  tags = {
    Name = "host-instance"
    App = var.AppName
  }
}

resource "aws_db_instance" "postgres-db" {
  allocated_storage = 20
  engine = "postgres"
  engine_version = "16"
  instance_class = "db.t3.micro"
  db_name = "movie_reviews"
  username = "movie_reviews"
  password = "CHANGE_ME"
  parameter_group_name = "default.postgres16"
  vpc_security_group_ids = [aws_security_group.allow-postgres.id]
  db_subnet_group_name = aws_db_subnet_group.main.name
  publicly_accessible = true
  identifier = "movie-reviews-db"
  allow_major_version_upgrade = true
  skip_final_snapshot = true

  tags = {
    Name = "movie-reviews-db"
    App = var.AppName
  }
}

resource "aws_ssm_parameter" "db_url" {
  name = "/movie-reviews/db-url"
  type = "SecureString"
  value = "postgres://${aws_db_instance.postgres-db.username}:${aws_db_instance.postgres-db.password}@${aws_db_instance.postgres-db.endpoint}/${aws_db_instance.postgres-db.db_name}"

  lifecycle {
    ignore_changes = [value]
  }
}

resource "aws_ssm_parameter" "secrets" {
  for_each = toset(var.secrets)
  name = each.value
  type = "SecureString"
  value = "CHANGE_ME"

  lifecycle {
    ignore_changes = [value]
  }
}

