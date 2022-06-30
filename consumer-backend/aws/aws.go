package aws

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// ToResponse transforms an interface into AWS response
func ToResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp.Body = string(stringBody)
	return &resp, nil
}
