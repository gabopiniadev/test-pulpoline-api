package ai

import "context"

type AIClient interface {
	ProcessText(ctx context.Context, text string) (string, error)
}
