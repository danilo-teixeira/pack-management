package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type (
	Client interface {
		Do(ctx context.Context, request Request, responseBody interface{}) error
	}

	Request struct {
		Method   string
		URL      string
		BodyJSON *string
	}

	client struct {
		stdClient *http.Client
	}
)

var (
	ErrRequestFailed = errors.New("request failed")
)

func NewClient() Client {
	return &client{
		stdClient: &http.Client{},
	}
}

func (c *client) Do(ctx context.Context, request Request, responseBody interface{}) error {
	var reqBody io.Reader
	if request.BodyJSON != nil {
		reqBody = bytes.NewBuffer([]byte(*request.BodyJSON))
	}

	if request.Method == "" {
		request.Method = http.MethodGet
	}

	req, err := http.NewRequestWithContext(ctx, request.Method, request.URL, reqBody)
	if err != nil {
		return err
	}

	response, err := c.stdClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return ErrRequestFailed
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}

	return nil
}
