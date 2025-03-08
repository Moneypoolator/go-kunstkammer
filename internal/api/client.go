package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type KaitenClient struct {
	client  *http.Client
	baseURL string
	token   string
}

func CreateKaitenClient(token string, kaitenURL string) *KaitenClient {
	//client := createHTTPClient()
	return &KaitenClient{
		client:  &http.Client{}, // lcient
		baseURL: kaitenURL,
		token:   token,
	}
}

func (kc *KaitenClient) doRequest(method, path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, kc.baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+kc.token)
	req.Header.Set("Content-Type", "application/json")

	return kc.client.Do(req)
}

func (kc *KaitenClient) doRequestWithBody(method, path string, body interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, kc.baseURL+path, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+kc.token)
	req.Header.Set("Content-Type", "application/json")

	return kc.client.Do(req)
}

func (kc *KaitenClient) decodeResponse(resp *http.Response, v interface{}) error {

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if v == nil {
		return nil // Если не нужно декодировать ответ
	}

	return json.NewDecoder(resp.Body).Decode(v)
}
