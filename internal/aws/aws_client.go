package aws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

const (
    DEFAULT_METHOD = "POST"
    SERVICE = "cloud9"
    URL_PATTERN = "https://%s.%s.amazonaws.com/"
    OPERATION_PREFIX = "AWSCloud9WorkspaceManagementService"
    AWS_JSON = "application/x-amz-json-1.1"
)

type AWSCloud9Client struct {
    region string
    client *http.Client
    service string
    signer *v4.Signer
    url string
}

func New(credentials *credentials.Credentials, region string) *AWSCloud9Client {
    return &AWSCloud9Client{
        region: region,
        client: &http.Client{},
        service: SERVICE,
        signer: v4.NewSigner(credentials),
        url: fmt.Sprintf(URL_PATTERN, SERVICE, region),
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
