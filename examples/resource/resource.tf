terraform {
  required_providers {
    awscloud9 = {
      source = "m1dugh/awscloud9"
    }
  }
}

provider "awscloud9" {
  region = var.region
}

variable "environment" {
  type = object({
    name     = string
    login    = string
    hostname = string
  })
  description = "The configuration for the environment"
}

variable "user_arn" {
  type        = string
  description = "The arn of the IAM role of user to give permissions to"
}

variable "region" {
  type = string
}

resource "awscloud9_ssh_environment" "test" {
  name       = var.environment.name
  login_name = var.environment.login
  hostname   = var.environment.hostname

  tags = {
    "managed-by" = "terraform"
    "owner"      = "infra"
  }
}

resource "awscloud9_environment_membership" "test" {
  environment_id = awscloud9_ssh_environment.test.environment_id
  permissions    = "read-write"
  user_arn       = var.user_arn
}

output "env_arn" {
  value = awscloud9_ssh_environment.test.arn
}
