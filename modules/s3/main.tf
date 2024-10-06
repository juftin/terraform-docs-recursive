/**
 * # s3
 *
 * This module creates a new S3 bucket and
 * applies some nice settings to it.
*/

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

variable "bucket_name" {
  type    = string
  default = "my-test-bucket"
}

provider "aws" {
  region = var.region
}

locals {
  bucket_name = var.bucket_name
}

resource "aws_s3_bucket" "bucket" {
  bucket = local.bucket_name

  tags = {
    Name        = local.bucket_name
    Environment = "dev"
  }
}

resource "aws_s3_bucket_public_access_block" "bucket_access" {
  bucket = aws_s3_bucket.bucket.bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_versioning" "bucket_versioning" {
  bucket = aws_s3_bucket.bucket.bucket
  versioning_configuration {
    status = "Enabled"
  }
}

output "name" {
  value       = aws_s3_bucket.bucket.bucket
  description = "Name of the bucket"
}

output "uri" {
  value       = "s3://${aws_s3_bucket.bucket.bucket}"
  description = "URI of the bucket"
}
