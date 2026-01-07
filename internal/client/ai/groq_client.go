package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"test-pulpoline-api/pkg/errors"
)

type GroqClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

type GroqRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type GroqResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

func NewGroqClient(apiKey string) *GroqClient {
	if apiKey == "" {
		log.Println("ADVERTENCIA: GROQ_API_KEY no está configurada")
	}

	return &GroqClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.groq.com/openai/v1/chat/completions",
	}
}

func (c *GroqClient) ProcessText(ctx context.Context, text string) (string, error) {
	if c.apiKey == "" {
		return "", errors.ErrMissingAPIKey
	}

	reqBody := GroqRequest{
		Model: "llama-3.1-8b-instant",
		Messages: []Message{
			{
				Role:    "user",
				Content: text,
			},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error al serializar la petición: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error al crear la petición: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al realizar la petición: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer la respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error de la API: %s - %s", resp.Status, string(body))
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", fmt.Errorf("error al parsear la respuesta: %w", err)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no se recibió ninguna respuesta de la API")
	}

	return groqResp.Choices[0].Message.Content, nil
}
