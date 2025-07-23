# AWS Integration Example
# This example shows how to use the YAML Flattener provider with AWS resources

terraform {
  required_version = ">= 1.0"
  required_providers {
    yamlflattener = {
      source  = "Perun-Engineering/yamlflattener"
      version = "~> 1.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "yamlflattener" {}

provider "aws" {
  region = local.aws_config["region"]
}

# Flatten AWS infrastructure configuration
data "yamlflattener_flatten" "aws_config" {
  yaml_content = <<-EOT
    infrastructure:
      aws:
        region: "us-west-2"
        availability_zones:
          - "us-west-2a"
          - "us-west-2b"
          - "us-west-2c"

        vpc:
          cidr: "10.0.0.0/16"
          enable_dns_hostnames: true
          enable_dns_support: true

        subnets:
          public:
            - cidr: "10.0.1.0/24"
              az: "us-west-2a"
            - cidr: "10.0.2.0/24"
              az: "us-west-2b"
          private:
            - cidr: "10.0.10.0/24"
              az: "us-west-2a"
            - cidr: "10.0.20.0/24"
              az: "us-west-2b"

        instances:
          web:
            type: "t3.medium"
            count: 2
            ami_filter: "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"
          database:
            engine: "postgres"
            version: "13.7"
            instance_class: "db.t3.micro"
            allocated_storage: 20

        security_groups:
          web:
            ingress:
              - from_port: 80
                to_port: 80
                protocol: "tcp"
                cidr_blocks: ["0.0.0.0/0"]
              - from_port: 443
                to_port: 443
                protocol: "tcp"
                cidr_blocks: ["0.0.0.0/0"]
          database:
            ingress:
              - from_port: 5432
                to_port: 5432
                protocol: "tcp"
                source_security_group: "web"

        tags:
          Environment: "production"
          Project: "yaml-flattener-demo"
          Owner: "devops-team"
  EOT
}

# Extract configuration into locals for easier access
locals {
  config = data.yamlflattener_flatten.aws_config.flattened

  # AWS-specific configuration
  aws_config = {
    region = local.config["infrastructure.aws.region"]
    vpc_cidr = local.config["infrastructure.aws.vpc.cidr"]
    enable_dns_hostnames = tobool(local.config["infrastructure.aws.vpc.enable_dns_hostnames"])
    enable_dns_support = tobool(local.config["infrastructure.aws.vpc.enable_dns_support"])
  }

  # Extract availability zones
  availability_zones = [
    local.config["infrastructure.aws.availability_zones[0]"],
    local.config["infrastructure.aws.availability_zones[1]"],
    local.config["infrastructure.aws.availability_zones[2]"]
  ]

  # Extract public subnets
  public_subnets = [
    {
      cidr = local.config["infrastructure.aws.subnets.public[0].cidr"]
      az   = local.config["infrastructure.aws.subnets.public[0].az"]
    },
    {
      cidr = local.config["infrastructure.aws.subnets.public[1].cidr"]
      az   = local.config["infrastructure.aws.subnets.public[1].az"]
    }
  ]

  # Extract private subnets
  private_subnets = [
    {
      cidr = local.config["infrastructure.aws.subnets.private[0].cidr"]
      az   = local.config["infrastructure.aws.subnets.private[0].az"]
    },
    {
      cidr = local.config["infrastructure.aws.subnets.private[1].cidr"]
      az   = local.config["infrastructure.aws.subnets.private[1].az"]
    }
  ]

  # Common tags
  common_tags = {
    Environment = local.config["infrastructure.aws.tags.Environment"]
    Project     = local.config["infrastructure.aws.tags.Project"]
    Owner       = local.config["infrastructure.aws.tags.Owner"]
  }
}

# Data source for AMI
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = [local.config["infrastructure.aws.instances.web.ami_filter"]]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = local.aws_config.vpc_cidr
  enable_dns_hostnames = local.aws_config.enable_dns_hostnames
  enable_dns_support   = local.aws_config.enable_dns_support

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-vpc"
  })
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-igw"
  })
}

# Public Subnets
resource "aws_subnet" "public" {
  count = length(local.public_subnets)

  vpc_id                  = aws_vpc.main.id
  cidr_block              = local.public_subnets[count.index].cidr
  availability_zone       = local.public_subnets[count.index].az
  map_public_ip_on_launch = true

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-public-${count.index + 1}"
    Type = "public"
  })
}

# Private Subnets
resource "aws_subnet" "private" {
  count = length(local.private_subnets)

  vpc_id            = aws_vpc.main.id
  cidr_block        = local.private_subnets[count.index].cidr
  availability_zone = local.private_subnets[count.index].az

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-private-${count.index + 1}"
    Type = "private"
  })
}

# Route Table for Public Subnets
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-public-rt"
  })
}

# Associate Public Subnets with Route Table
resource "aws_route_table_association" "public" {
  count = length(aws_subnet.public)

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Security Group for Web Servers
resource "aws_security_group" "web" {
  name_prefix = "${local.common_tags.Project}-web-"
  vpc_id      = aws_vpc.main.id

  # HTTP
  ingress {
    from_port   = tonumber(local.config["infrastructure.aws.security_groups.web.ingress[0].from_port"])
    to_port     = tonumber(local.config["infrastructure.aws.security_groups.web.ingress[0].to_port"])
    protocol    = local.config["infrastructure.aws.security_groups.web.ingress[0].protocol"]
    cidr_blocks = [local.config["infrastructure.aws.security_groups.web.ingress[0].cidr_blocks[0]"]]
  }

  # HTTPS
  ingress {
    from_port   = tonumber(local.config["infrastructure.aws.security_groups.web.ingress[1].from_port"])
    to_port     = tonumber(local.config["infrastructure.aws.security_groups.web.ingress[1].to_port"])
    protocol    = local.config["infrastructure.aws.security_groups.web.ingress[1].protocol"]
    cidr_blocks = [local.config["infrastructure.aws.security_groups.web.ingress[1].cidr_blocks[0]"]]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-web-sg"
  })
}

# Security Group for Database
resource "aws_security_group" "database" {
  name_prefix = "${local.common_tags.Project}-db-"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = tonumber(local.config["infrastructure.aws.security_groups.database.ingress[0].from_port"])
    to_port         = tonumber(local.config["infrastructure.aws.security_groups.database.ingress[0].to_port"])
    protocol        = local.config["infrastructure.aws.security_groups.database.ingress[0].protocol"]
    security_groups = [aws_security_group.web.id]
  }

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-db-sg"
  })
}

# Web Server Instances
resource "aws_instance" "web" {
  count = tonumber(local.config["infrastructure.aws.instances.web.count"])

  ami                    = data.aws_ami.ubuntu.id
  instance_type          = local.config["infrastructure.aws.instances.web.type"]
  subnet_id              = aws_subnet.public[count.index % length(aws_subnet.public)].id
  vpc_security_group_ids = [aws_security_group.web.id]

  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl start nginx
    systemctl enable nginx
    echo "<h1>Web Server ${count.index + 1}</h1>" > /var/www/html/index.html
  EOF

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-web-${count.index + 1}"
    Type = "web-server"
  })
}

# Database Subnet Group
resource "aws_db_subnet_group" "main" {
  name       = "${local.common_tags.Project}-db-subnet-group"
  subnet_ids = aws_subnet.private[*].id

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-db-subnet-group"
  })
}

# RDS Database Instance
resource "aws_db_instance" "main" {
  identifier = "${local.common_tags.Project}-database"

  engine         = local.config["infrastructure.aws.instances.database.engine"]
  engine_version = local.config["infrastructure.aws.instances.database.version"]
  instance_class = local.config["infrastructure.aws.instances.database.instance_class"]

  allocated_storage = tonumber(local.config["infrastructure.aws.instances.database.allocated_storage"])
  storage_type      = "gp2"

  db_name  = "appdb"
  username = "admin"
  password = "changeme123!" # In production, use AWS Secrets Manager

  vpc_security_group_ids = [aws_security_group.database.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name

  skip_final_snapshot = true

  tags = merge(local.common_tags, {
    Name = "${local.common_tags.Project}-database"
  })
}

# Application Load Balancer
resource "aws_lb" "main" {
  name               = "${local.common_tags.Project}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.web.id]
  subnets            = aws_subnet.public[*].id

  tags = local.common_tags
}

# ALB Target Group
resource "aws_lb_target_group" "web" {
  name     = "${local.common_tags.Project}-web-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    path                = "/"
    matcher             = "200"
  }

  tags = local.common_tags
}

# ALB Target Group Attachments
resource "aws_lb_target_group_attachment" "web" {
  count = length(aws_instance.web)

  target_group_arn = aws_lb_target_group.web.arn
  target_id        = aws_instance.web[count.index].id
  port             = 80
}

# ALB Listener
resource "aws_lb_listener" "web" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.web.arn
  }
}

# Outputs
output "load_balancer_dns" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "database_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "web_instance_ips" {
  description = "Public IP addresses of web instances"
  value       = aws_instance.web[*].public_ip
}

output "flattened_config_sample" {
  description = "Sample of flattened configuration keys"
  value = {
    vpc_cidr        = local.config["infrastructure.aws.vpc.cidr"]
    instance_type   = local.config["infrastructure.aws.instances.web.type"]
    database_engine = local.config["infrastructure.aws.instances.database.engine"]
    environment     = local.config["infrastructure.aws.tags.Environment"]
  }
}
