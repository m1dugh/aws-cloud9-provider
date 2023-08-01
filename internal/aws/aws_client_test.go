package aws

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

var access_key_id = os.Getenv("AWS_ACCESS_KEY_ID")
var secret_access_key = os.Getenv("AWS_SECRET_ACCESS_KEY")

const region = "eu-west-3"

func TestRequest(t *testing.T) {

    client := New(context.Background(), credentials.NewStaticCredentials(access_key_id, secret_access_key, ""), region)
    // environment, err := client.DescribeSSHRemote("573a64362bc44311a52fa6e0178b3dd3")
    envs, err := client.GetSSHEnvironments([]string{"573a64362bc44311a52fa6e0178b3dd3"})
    if err != nil {
        t.Fatalf("error: %s", err)
    }

    env := envs[0]
    t.Fatalf("found name: %s => %s@%s", env.Name, env.LoginName, env.Hostname)
}
