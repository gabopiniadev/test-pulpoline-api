package service

import (
	"context"
	"fmt"

	"test-pulpoline-api/internal/client/ai"
	"test-pulpoline-api/internal/queue/models"
)

type Service struct {
	aiClient ai.AIClient
}

func NewService(aiClient ai.AIClient) *Service {
	return &Service{
		aiClient: aiClient,
	}
}

func (s *Service) ProcessText(ctx context.Context, text string) (string, error) {
	return s.aiClient.ProcessText(ctx, text)
}

func (s *Service) ProcessRequest(
	ctx context.Context,
	requestID string,
	text string,
	resultChan chan<- models.Response,
	errorChan chan<- error,
) {
	response, err := s.aiClient.ProcessText(ctx, text)
	if err != nil {
		errorChan <- fmt.Errorf("error al procesar con IA: %w", err)
		return
	}

	resultChan <- models.Response{
		ID:       requestID,
		Text:     text,
		Response: response,
	}
}
