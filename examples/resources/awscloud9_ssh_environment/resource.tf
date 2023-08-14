# Basic ssh environment
resource "awscloud9_ssh_environment" "env" {
  name       = "my_environment"
  login_name = "my_user"
  hostname   = "my-host.ec2.amazonaws.com"

  tags = {
    "managed-by" = "terraform"
    "owner"      = "infra"
  }
}

# SSH Environment with bastion_url
resource "awscloud9_ssh_environment" "env_with_bastion" {
  name       = "my_protected_environment"
  login_name = "my_user"
  hostname   = "my-host.ec2.amazonaws.com"

  description = "An SSH environment for my_user@my-host"

  port             = 42
  environment_path = "/tmp/folders/my_user"
  node_path        = "/bin/node"

  bastion_url = "my_user@my.proxy.com:22"
}
