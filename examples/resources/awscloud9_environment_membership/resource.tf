data "awscloud9_ssh_environment" "env" {
  environment_id = "..."
}

resource "awscloud9_environment_membership" "membership" {
  environment_id = awscloud9_ssh_environment.env.environment_id
  permissions    = "read-write"
  user_arn       = "arn:aws:..."
}
