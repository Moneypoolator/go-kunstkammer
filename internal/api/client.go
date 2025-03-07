package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://kaiten.norsoft.ru/api/latest"
)

type KaitenClient struct {
	client  *http.Client
	baseURL string
	token   string
}

// func createHTTPClient() *http.Client {
// 	return &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: true, // Отключаем проверку сертификата
// 			},
// 		},
// 	}
// }

// func createTLSConfig(caCertPath string) (*tls.Config, error) {
// 	caCert, err := os.ReadFile(caCertPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	return &tls.Config{
// 		RootCAs: caCertPool,
// 	}, nil
// }

// func createHTTPClientWithCertificate(crtFileName string) *http.Client {
// 	tlsConfig, err := createTLSConfig(crtFileName)
// 	if err != nil {
// 		log.Fatalf("Failed to create TLS config: %v", err)
// 	}

// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: tlsConfig,
// 		},
// 	}
// 	return client
// }

func CreateKaitenClient(token string) *KaitenClient {
	//client := createHTTPClient()
	return &KaitenClient{
		client:  &http.Client{}, // lcient
		baseURL: baseURL,
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
	return json.NewDecoder(resp.Body).Decode(v)
}
