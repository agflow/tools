package accountmanagement

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/agflow/tools/consumer-backend/model"
	"github.com/agflow/tools/log"
)

type AccountManagement struct {
	URL string
	API
}

type API interface {
	GetUser(string) (*model.User, error)
}

func New(url string) *AccountManagement {
	return &AccountManagement{
		URL: url,
	}
}

func (am *AccountManagement) GetUser(token string) (
	*model.User, error,
) {
	client := http.Client{}
	req, err := http.NewRequest("GET", am.URL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("can't create account management request, err: %w", err)
	}
	req.Header = http.Header{
		"fusionauth": []string{token},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request to account management, err: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("account management code %d, response: %v", resp.StatusCode, resp)
	}
	defer func() { log.IfErrorDiffNil(resp.Body.Close()) }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read account management response, err: %w", err)
	}

	user := model.User{}
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("can't unmarshal user, err: %w", err)
	}

	return &user, nil
}
