
# Configuration-based authentication
provider "awscloud9" {
    aws_access_key_id = "..."
    aws_secret_access_key = "..."
    region = "us-east-1"
}

# Environment variables-based authentication
provider "awscloud9" {
    region = "us-east-1"
}
