# Terraform Plan (Minimal)

This folder is a placeholder for infrastructure as code.

## Minimal target architecture

- VPC with public/private subnets
- RDS (PostgreSQL + PostGIS) in private subnets
- ECS/Fargate or EC2 for the Go API
- ALB for public ingress

## Suggested next steps

1. Create `main.tf` with AWS provider and a VPC module.
2. Add an RDS module (PostgreSQL 15 with PostGIS).
3. Add ECS or EC2 module for the API.
4. Add networking outputs for service discovery.

When you want, I can scaffold the real Terraform modules and variables.
