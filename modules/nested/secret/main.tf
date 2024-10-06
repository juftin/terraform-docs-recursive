/**
 * # secret
 *
 * This module creates a secret in AWS Secrets Manager.
 * This is also a nested module, which will still be found.
*/

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "secret_name" {
  type    = string
  default = "my-test-bucket"
}

variable "secret_value" {
  type      = string
  default   = "my-test-bucket"
  sensitive = true
}

provider "aws" {
  region = var.region
}

locals {
  secret_name  = var.secret_name
  secret_value = var.secret_value
}

resource "aws_secretsmanager_secret" "secret" {
  name = local.secret_name
}

resource "aws_secretsmanager_secret_version" "secret_version" {
  secret_id     = aws_secretsmanager_secret.secret.id
  secret_string = local.secret_value
}
