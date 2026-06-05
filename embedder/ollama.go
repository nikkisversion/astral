package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	client     *http.Client
	baseURL    string
	embedModel string
	chatModel  string
}

// NewOllamaEmbedder initializes the client.
// Default Ollama URL is usually http://localhost:11434
func NewOllamaEmbedder(baseURL, embedModel, chatModel string) *OllamaEmbedder {
	return &OllamaEmbedder{
		client:     &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		embedModel: embedModel,
		chatModel:  chatModel,
	}
}

// EmbedStrings takes a batch of text strings and returns their vector representations.
func (oe *OllamaEmbedder) EmbedStrings(ctx context.Context, inputs []string) ([][]float32, error) {

	reqBody := OllamaEmbedRequest{
		Model: oe.embedModel,
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

func (oe *OllamaEmbedder) GenerateAnswer(ctx context.Context, messages []ChatMessage) (string, error) {

	payload := map[string]any{
		"model":    oe.chatModel,
		"messages": messages,
		"stream":   false,
	}

	body, errBody := json.Marshal(payload)
	if errBody != nil {
		return "", errors.New("Failed to encode request body: " + errBody.Error())
	}

	req, errReq := http.NewRequestWithContext(ctx, http.MethodPost, oe.baseURL+"/api/chat", bytes.NewReader(body))
	if errReq != nil {
		return "", errors.New("Failed to create request: " + errReq.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	resp, errResp := oe.client.Do(req)
	if errResp != nil {
		return "", errors.New("Failed to send request: " + errResp.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Unexpected status code: " + resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read body: %w", err)
	}

	var res ChatResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return res.Message.Content, nil

}
