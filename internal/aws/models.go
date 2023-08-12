package aws

const (
    OWNER = "owner"
    READ_WRITE = "read-write"
    READONLY = "read-only"
)

type Cloud9EnvironmentMembership struct {
    EnvironmentId string `json:"environment_id"`
    Permissions string `json:"permissions"`
    UserARN string `json:"userArn"`
    UserID string `json:"userId"`
}

type Cloud9SSHEnvironment struct {
    EnvironmentId string `json:"environment_id"`
    Name string `json:"name"` 
    Description string `json:"description,omitempty"`
    LoginName string `json:"loginName"`
    Hostname string `json:"host"`
    Port int16 `json:"port"`
    EnvironmentPath string `json:"environmentPath,omitempty"`
    NodePath string `json:"nodePath,omitempty"`
    BastionHost string `json:"bastionHost,omitempty"`
    DryRun bool `json:"dryRun"`
    Memberships []Cloud9EnvironmentMembership `json:"memberships"`
    Tags map[string]string `json:"tags"`
}

type CreateEnvironmentSSHRequest struct {
    Name string `json:"name"` 
    Description string `json:"description,omitempty"`
    LoginName string `json:"loginName"`
    Hostname string `json:"host"`
    Port int16 `json:"port"`
    EnvironmentPath string `json:"environmentPath,omitempty"`
    NodePath string `json:"nodePath,omitempty"`
    BastionHost string `json:"bastionHost,omitempty"`
    DryRun bool `json:"dryRun"`
    Tags []string `json:"tags"`
}

type CreateEnvironmentSSHResult struct {
    EnvironmentId string `json:"environmentId"`
}

type DescribeSSHRemoteRequest struct {
    EnvironmentId string `json:"environmentId"`
}

type SSHRemoteEnvironmentDescription struct {
    EnvironmentPath string `json:"environmentPath"`
    Hostname string `json:"host"`
    Description string `json:"description,omitempty"`
    LoginName string `json:"loginName"`
    Port int16 `json:"port"`
    NodePath string `json:"nodePath"`
    BastionHost string `json:"bastionHost"`
}

type DescribeSSHRemoteResult struct {
    Results SSHRemoteEnvironmentDescription `json:"remote"`
}

type UpdateSSHRemoteRequest struct {
    EnvironmentId string `json:"environmentId"`
    Name string `json:"name"` 
    Description string `json:"description,omitempty"`
    LoginName string `json:"loginName"`
    Hostname string `json:"host"`
    Port int16 `json:"port"`
    EnvironmentPath string `json:"environmentPath"`
    NodePath string `json:"nodePath"`
    BastionHost string `json:"bastionHost"`
}

type GetUserPublicKeyResult struct {
    PublicKey string `json:"publicKey"`
}

type AWSError struct {
    ExceptionType string `json:"__type"`
    Message string `json:"message"`
}
