package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type OllamaEmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type OllamaEmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

type OllamaEmbedder struct {
	client  *http.Client
	baseURL string
	model   string
}

// NewOllamaEmbedder initializes the client.
// Default Ollama URL is usually http://localhost:11434
func NewOllamaEmbedder(baseURL, model string) *OllamaEmbedder {
	return &OllamaEmbedder{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: baseURL,
		model:   model,
	}
}

// EmbedStrings takes a batch of text strings and returns their vector representations.
func (oe *OllamaEmbedder) EmbedStrings(ctx context.Context, inputs []string) ([][]float32, error) {

	reqBody := OllamaEmbedRequest{
		Model: oe.model,
		Input: inputs,
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return nil, errors.New("Failed to encode request body: " + err.Error())
	}

	url := oe.baseURL + "/api/embed"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, errors.New("Failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oe.client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to send request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unexpected status code: " + resp.Status)
	}

	var embedResp OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, errors.New("Failed to decode response: " + err.Error())
	}

	return embedResp.Embeddings, nil
}
