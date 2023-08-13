package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/cloud9"
)

const (
    DEFAULT_METHOD = "POST"
    SERVICE = "cloud9"
    URL_PATTERN = "https://%s.%s.amazonaws.com/"
    OPERATION_PREFIX = "AWSCloud9WorkspaceManagementService"
    AWS_JSON = "application/x-amz-json-1.1"
    MAX_RESULTS = 25
)

type AWSCloud9Client struct {
    region string
    client *http.Client
    service string
    signer *v4.Signer
    url string
    Cloud9 *cloud9.Cloud9
    ctx context.Context
    session *session.Session
}

func New(ctx context.Context, credentials *credentials.Credentials, region string) *AWSCloud9Client {
    config := aws.NewConfig().WithRegion(region)
    session := session.Must(session.NewSession())

    client := cloud9.New(session, config)
    return &AWSCloud9Client{
        region: region,
        client: &http.Client{},
        service: SERVICE,
        signer: v4.NewSigner(credentials),
        url: fmt.Sprintf(URL_PATTERN, SERVICE, region),
        ctx: ctx,
        Cloud9: client,
    }
}

func (client *AWSCloud9Client) signRequest(request *http.Request, body io.ReadSeeker) error {

    _, err := client.signer.Sign(request, body, client.service, client.region, time.Now())
    if err != nil {
        return err
    }
    return nil
}

func (client *AWSCloud9Client) executeCloud9(operation string, body interface{}) (*http.Response, error) {

    bodyString, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }

    request, err := http.NewRequest(DEFAULT_METHOD, client.url, bytes.NewReader(bodyString))
    if err != nil {
        return nil, err
    }
    opString := fmt.Sprintf("%s.%s", OPERATION_PREFIX, operation)

    request.Header.Set("Content-Type", AWS_JSON)
    request.Header.Set("X-Amz-Target", opString)

    if err = client.signRequest(request, bytes.NewReader(bodyString)); err != nil {
        return nil, err
    }

    return client.client.Do(request)
}

func (client *AWSCloud9Client) GetUserPublicKey() (*GetUserPublicKeyResult, error) {
    var body struct{}
    res, err := client.executeCloud9("GetUserPublicKey", body)
    if err != nil {
        return nil, err
    }

    bodyBytes, err := io.ReadAll(res.Body)
    defer res.Body.Close()
    if err != nil {
        return nil, err
    }
    if res.StatusCode == 400 {
        var error AWSError
        if err = json.Unmarshal(bodyBytes, &error); err != nil {
            return nil, err
        }
        return nil, errors.New(fmt.Sprintf("An error occured of type %s: %s", error.ExceptionType, error.Message))
    }
    var result GetUserPublicKeyResult
    if err = json.Unmarshal(bodyBytes, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func (client *AWSCloud9Client) DescribeSSHRemote(environmentId string) (*DescribeSSHRemoteResult, error) {
    request := DescribeSSHRemoteRequest{
        EnvironmentId: environmentId,
    }

    res, err := client.executeCloud9("DescribeSSHRemote", request)
    if err != nil {
        return nil, err
    }
    bodyBytes, err := io.ReadAll(res.Body)
    defer res.Body.Close()
    if err != nil {
        return nil, err
    }

    if res.StatusCode == 400 {
        var error AWSError
        if err = json.Unmarshal(bodyBytes, &error); err != nil {
            return nil, err
        }
        return nil, errors.New(fmt.Sprintf("An error occured of type %s: %s", error.ExceptionType, error.Message))
    }
    var result DescribeSSHRemoteResult
    if err = json.Unmarshal(bodyBytes, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func (client *AWSCloud9Client) UpdateSSHRemote(request *UpdateSSHRemoteRequest) error {

    res, err := client.executeCloud9("UpdateSSHRemote", request)
    if err != nil {
        return err
    }
    bodyBytes, err := io.ReadAll(res.Body)
    defer res.Body.Close()
    if err != nil {
        return err
    }

    if res.StatusCode == 400 {
        var error AWSError
        if err = json.Unmarshal(bodyBytes, &error); err != nil {
            return err
        }
        return errors.New(fmt.Sprintf("An error occured of type %s: %s", error.ExceptionType, error.Message))
    }

    return nil
}

func (client *AWSCloud9Client) CreateEnvironmentSSH(request *CreateEnvironmentSSHRequest) (*CreateEnvironmentSSHResult, error) {
    res, err := client.executeCloud9("CreateEnvironmentSSH", request)
    if err != nil {
        return nil, err
    }
    bodyBytes, err := io.ReadAll(res.Body)
    defer res.Body.Close()
    if err != nil {
        return nil, err
    }

    if res.StatusCode == 400 {
        var error AWSError
        if err = json.Unmarshal(bodyBytes, &error); err != nil {
            return nil, err
        }
        return nil, errors.New(fmt.Sprintf("An error occured of type %s: %s", error.ExceptionType, error.Message))
    }
    var result CreateEnvironmentSSHResult
    if err = json.Unmarshal(bodyBytes, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func (client *AWSCloud9Client) GetMemberShips(environmentId string) ([]Cloud9EnvironmentMembership, error) {


    input := &cloud9.DescribeEnvironmentMembershipsInput {
        EnvironmentId: aws.String(environmentId),
    }

    var res []Cloud9EnvironmentMembership = make([]Cloud9EnvironmentMembership, 0)

    var hasResults bool = true
    for hasResults {
        response, err := client.Cloud9.DescribeEnvironmentMemberships(input)
        if err != nil {
            return nil, err
        }

        for _, membership := range response.Memberships {
            res = append(res, Cloud9EnvironmentMembership{
                EnvironmentId: *membership.EnvironmentId,
                Permissions: *membership.Permissions,
                UserARN: *membership.UserArn,
                UserID: *membership.UserId,
            })
        }

        if response.NextToken != nil {
            input.NextToken = response.NextToken
        } else {
            hasResults = false
        }
    }

    return res, nil
}

func (client *AWSCloud9Client) GetSSHEnvironments(envIds ...string) ([]Cloud9SSHEnvironment, error) {
    var res []Cloud9SSHEnvironment = make([]Cloud9SSHEnvironment, 0, len(envIds))
    cursor := 0
    for ;cursor < len(envIds); cursor += MAX_RESULTS {
        var ids []*string
        if cursor + MAX_RESULTS < len(envIds) {
            ids = make([]*string, MAX_RESULTS)
        } else {
            ids = make([]*string, len(envIds) - cursor)
        }

        for i := 0; i < len(ids); i++ {
            ids[i] = &envIds[cursor + i]
        }

        response, err := client.Cloud9.DescribeEnvironments(&cloud9.DescribeEnvironmentsInput{
            EnvironmentIds: ids,
        })

        if err != nil {
            return nil, err
        }

        for _, env := range response.Environments {

            envId := *env.Id
            sshConfig, err := client.DescribeSSHRemote(envId)
            if err != nil {
                return nil, err
            }

            tags, err := client.Cloud9.ListTagsForResource(&cloud9.ListTagsForResourceInput{
                ResourceARN: env.Arn,
            })

            if err != nil {
                return nil, err
            }

            tagMap := make([]Tag, 0)
            for _, tag := range tags.Tags {
                tagMap = append(tagMap, Tag{
                    Key: *tag.Key,
                    Value: *tag.Value,
                })
            }

            res = append(res, Cloud9SSHEnvironment{
                Arn: *env.Arn,
                EnvironmentId: envId,
                Name: *env.Name,
                Description: *env.Description,
                EnvironmentPath: sshConfig.Results.EnvironmentPath,
                Hostname: sshConfig.Results.Hostname,
                LoginName: sshConfig.Results.LoginName,
                Port: sshConfig.Results.Port,
                NodePath: sshConfig.Results.NodePath,
                BastionHost: sshConfig.Results.BastionHost,
                Tags: tagMap,
            })
        }
    }

    return res, nil
}

func (client *AWSCloud9Client) UpdateEnvironment(env Cloud9SSHEnvironment) error {
    _, err := client.Cloud9.UpdateEnvironment(&cloud9.UpdateEnvironmentInput{
        EnvironmentId: &env.EnvironmentId,
        Name: &env.Name,
        Description: &env.Description,
    })

    if err != nil {
        return err
    }

    var updateRequest UpdateSSHRemoteRequest
    updateRequest.EnvironmentId = env.EnvironmentId
    updateRequest.Hostname = env.Hostname
    updateRequest.LoginName = env.LoginName
    updateRequest.NodePath = env.NodePath
    updateRequest.BastionHost = env.BastionHost
    updateRequest.Port = env.Port
    updateRequest.EnvironmentPath = env.EnvironmentPath

    err = client.UpdateSSHRemote(&updateRequest)
    if err != nil {
        return err
    }

    return nil
}
