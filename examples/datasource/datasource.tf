terraform {
  required_providers {
    awscloud9 = {
      source = "m1dugh/awscloud9"
    }
  }
}

variable "env_id" {
  type        = string
  description = "An environment id"
}

provider "awscloud9" {
  region = "eu-west-3"
}


data "awscloud9_ssh_environment" "test" {
  environment_id = var.env_id
}

output "test" {
  value = data.awscloud9_ssh_environment.test
}
