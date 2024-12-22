package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"io"
	"net/http"
	"time"
)

type Client interface {
	Do(ctx context.Context, method, url string, body []byte) ([]byte, error)
}

type APIConfig struct {
	Endpoint string        `envconfig:"APP_TODO_SERVICE_ENDPOINT" default:"http://localhost:5000"`
	Timeout  time.Duration `envconfig:"APP_TODO_SERVICE_TIMEOUT" default:"10s"`
	Port     string        `envconfig:"APP_TODO_SERVICE_PORT" default:"8080"`
}

type client struct {
	httpClient *http.Client
	apiConfig  APIConfig
}

func NewTodoServiceClient(httpClient *http.Client, apiConfig APIConfig) Client {

	return &client{
		httpClient: httpClient,
		apiConfig:  apiConfig,
	}
}

func (c *client) Do(ctx context.Context, method, url string, body []byte) ([]byte, error) {
	log.C(ctx).Infof("do in client for url %s and method %s", url, method)
	endpoint := c.apiConfig.Endpoint + url
	requestBody := []byte("{}")
	if body != nil {
		requestBody = body
	}

	log.C(ctx).Debugf("request body: %s", string(requestBody))

	token, ok := ctx.Value(constants.TokenCtxKey).(string)
	if !ok || token == "" {
		log.C(ctx).Error("authorization token missing or invalid in context")
		return nil, fmt.Errorf("authorization token missing or invalid")
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		log.C(ctx).Errorf("error creating request for url %s and method %s", endpoint, method)
		return nil, err
	}
	req.Header.Set("Content-Type", constants.ContentTypeJSON)
	req.Header.Set(constants.AuthorizationHeader, fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.C(ctx).Errorf("error doing request for url %s and method %s", endpoint, method)
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	log.C(ctx).Debugf("response status code: %d", resp.StatusCode)

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.C(ctx).Errorf("error closing response body for url %s and method %s", endpoint, method)
			return
		}
	}()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.C(ctx).Errorf("error reading response body for url %s and method %s", endpoint, method)
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}
