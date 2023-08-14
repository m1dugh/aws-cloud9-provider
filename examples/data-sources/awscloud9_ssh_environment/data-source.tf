data "awscloud9_ssh_environment" "env" {
  environment_id = var.env_id
}

output "env_arn" {
  value = data.awscloud9_ssh_environment.env.arn
}
