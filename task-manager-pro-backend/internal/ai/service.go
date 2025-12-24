package ai

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{provider: provider}
}

func (s *Service) SuggestTitles(ctx context.Context, description string) ([]string, error) {
	desc := strings.TrimSpace(description)
	if len(desc) < 10 || len(desc) > 300 {
		return nil, errors.New("invalid description")
	}

	titles, err := s.provider.SuggestTitles(ctx, desc)
	if err != nil {
		return nil, err
	}

	// Normalize + enforce constraints (defensivo)
	clean := make([]string, 0, 3)
	seen := map[string]struct{}{}

	for _, t := range titles {
		tt := strings.TrimSpace(t)
		tt = strings.TrimSuffix(tt, ".") // sem ponto final
		if tt == "" {
			continue
		}
		if len(tt) > 60 {
			tt = tt[:60]
		}
		key := strings.ToLower(tt)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		clean = append(clean, tt)
		if len(clean) == 3 {
			break
		}
	}

	// Se o provider voltar ruim, faz fallback mínimo
	for len(clean) < 3 {
		clean = append(clean, "Revisar tarefa")
	}

	return clean, nil
}

func (s *Service) ImproveDescription(ctx context.Context, title, description string) (string, []string, error) {
	t := strings.TrimSpace(title)
	d := strings.TrimSpace(description)

	if len(t) < 3 || len(t) > 80 || len(d) < 10 || len(d) > 1000 {
		return "", nil, errors.New("invalid input")
	}

	improved, bullets, err := s.provider.ImproveDescription(ctx, t, d)
	if err != nil {
		return "", nil, err
	}

	improved = strings.TrimSpace(improved)
	if improved == "" {
		improved = d // fallback: mantém o original
	}

	// normaliza bullets
	clean := make([]string, 0, 8)
	for _, b := range bullets {
		bb := strings.TrimSpace(b)
		if bb == "" {
			continue
		}
		if len(bb) > 120 {
			bb = bb[:120]
		}
		clean = append(clean, bb)
		if len(clean) == 8 {
			break
		}
	}

	// se vier vazio, cria bullets simples
	if len(clean) == 0 {
		clean = []string{"Definir passos", "Executar", "Testar"}
	}

	return improved, clean, nil
}
