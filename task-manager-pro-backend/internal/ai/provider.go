package ai

import "context"

type Provider interface {
	SuggestTitles(ctx context.Context, description string) ([]string, error)
	ImproveDescription(ctx context.Context, title, description string) (string, []string, error)
	GenerateSubtasks(ctx context.Context, title, description string) ([]SubtaskDTO, error)
}
