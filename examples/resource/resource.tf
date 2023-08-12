terraform {
  required_providers {
    awscloud9 = {
        source = "m1dugh/awscloud9"
    }
  }
}

provider "awscloud9" {
    region = "eu-west-3"
}

variable "environment" {
    type = object({
      name = string
      login = string
      hostname = string
    })
}

resource "awscloud9_ssh_environment" "test" {
    name = var.environment.name
    login_name = var.environment.login
    hostname = var.environment.hostname
}
