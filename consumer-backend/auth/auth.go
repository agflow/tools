package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/agflow/tools/consumer-backend/model"
)

type AWSLambdaResponse struct {
	StatusCode        int               `json:"statusCode"`
	Headers           map[string]string `json:"headers"`
	MultiValueHeaders map[string]string `json:"multiValueHeaders"`
	Body              string            `json:"body"`
}

// GetUser gets user from authentication system
func GetUser(token string) (*model.User, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("eu-west-1")})
	payload := struct {
		Headers map[string]string `json:"headers"`
	}{
		Headers: map[string]string{
			"Fusionauth": token,
		},
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	result, err := client.Invoke(
		&lambda.InvokeInput{
			FunctionName:   aws.String("auth-AuthFunction-C5c7YBM7DRGD"),
			InvocationType: aws.String("RequestResponse"),
			Payload:        payloadByte,
		},
	)
	if err != nil || *result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", "Error calling AuthFunction", err)
	}

	var resp AWSLambdaResponse
	if err := json.Unmarshal(result.Payload, &resp); err != nil {
		return nil, err
	}
	var user model.User
	if err := json.Unmarshal([]byte(resp.Body), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
