package ai

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

type OpenAIProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAIProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 25 * time.Second},
	}
}

func (p *OpenAIProvider) SuggestTitles(ctx context.Context, description string) ([]string, error) {
	prompt := fmt.Sprintf(`Você é um assistente que gera títulos curtos e claros para tarefas.

Gere exatamente 3 títulos curtos, objetivos e acionáveis com base na descrição abaixo.

Regras:
- Cada título deve ter no máximo 60 caracteres
- Não use emojis
- Não use ponto final
- Não numere
- Não explique nada
- Retorne APENAS um JSON válido no formato abaixo

Formato obrigatório:
{"titles":["","",""]}

Descrição da tarefa:
%s`, description)

	reqBody := map[string]any{
		"model": p.model,
		"input": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "input_text", "text": prompt},
				},
			},
		},
		// Força JSON
		"text": map[string]any{
			"format": map[string]any{
				"type": "json_object",
			},
		},
	}

	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		// Retorna body pra facilitar debug
		return nil, fmt.Errorf("openai error: status=%d body=%s", resp.StatusCode, string(rawBytes))
	}

	// 1) tenta extrair "output_text"
	// 2) se vazio, tenta extrair de output[].content[].text
	modelText, err := extractResponsesText(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract model text: %w | raw=%s", err, string(rawBytes))
	}
	if modelText == "" {
		return nil, fmt.Errorf("openai returned empty text | raw=%s", string(rawBytes))
	}

	// Esperamos JSON {"titles":[...]}
	var out struct {
		Titles []string `json:"titles"`
	}
	if err := json.Unmarshal([]byte(modelText), &out); err != nil {
		return nil, fmt.Errorf("failed to parse model json: %w | model_text=%s", err, modelText)
	}
	if len(out.Titles) == 0 {
		return nil, errors.New("no titles returned")
	}

	return out.Titles, nil
}

// extractResponsesText tenta pegar o texto do retorno da Responses API.
// Prioridade: output_text -> output[].content[].text
func extractResponsesText(raw []byte) (string, error) {
	// Estrutura mínima e tolerante
	var r struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := json.Unmarshal(raw, &r); err != nil {
		return "", err
	}

	if r.OutputText != "" {
		return r.OutputText, nil
	}

	for _, o := range r.Output {
		for _, c := range o.Content {
			if c.Text != "" {
				return c.Text, nil
			}
		}
	}

	return "", nil
}

func (p *OpenAIProvider) ImproveDescription(ctx context.Context, title, description string) (string, []string, error) {
	prompt := fmt.Sprintf(`Você é um assistente que melhora descrições de tarefas.

Com base no TÍTULO e na DESCRIÇÃO, reescreva a descrição deixando-a:
- clara, objetiva e acionável
- em português do Brasil
- sem inventar requisitos absurdos
- sem emojis

Retorne APENAS um JSON válido no formato:

{
  "improved_description": "texto...",
  "bullets": ["item1", "item2", "item3", "item4", "item5"]
}

Regras:
- improved_description: 2 a 6 linhas, sem markdown
- bullets: 3 a 6 itens, curtos, começando com verbo (ex: "Criar", "Testar", "Validar")

TÍTULO:
%s

DESCRIÇÃO:
%s`, title, description)

	reqBody := map[string]any{
		"model": p.model,
		"input": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "input_text", "text": prompt},
				},
			},
		},
		"text": map[string]any{
			"format": map[string]any{
				"type": "json_object",
			},
		},
	}

	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(b))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode >= 400 {
		return "", nil, fmt.Errorf("openai error: status=%d body=%s", resp.StatusCode, string(rawBytes))
	}

	modelText, err := extractResponsesText(rawBytes)
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract model text: %w | raw=%s", err, string(rawBytes))
	}
	if modelText == "" {
		return "", nil, fmt.Errorf("openai returned empty text | raw=%s", string(rawBytes))
	}

	var out struct {
		ImprovedDescription string   `json:"improved_description"`
		Bullets             []string `json:"bullets"`
	}

	if err := json.Unmarshal([]byte(modelText), &out); err != nil {
		return "", nil, fmt.Errorf("failed to parse model json: %w | model_text=%s", err, modelText)
	}

	return out.ImprovedDescription, out.Bullets, nil
}

func (p *OpenAIProvider) GenerateSubtasks(ctx context.Context, title, description string) ([]SubtaskDTO, error) {
	prompt := fmt.Sprintf(`Você é um assistente que cria um checklist de subtarefas.

Com base no TÍTULO e na DESCRIÇÃO, gere subtarefas:
- objetivas e acionáveis
- em português do Brasil
- começando com verbo (ex: "Criar", "Testar", "Validar")
- sem emojis
- sem itens genéricos tipo "Fazer", "Resolver", "Organizar"

Retorne APENAS um JSON válido no formato:

{
  "subtasks": [
    {"title": "texto", "done": false},
    {"title": "texto", "done": false}
  ]
}

Regras:
- Gere de 5 a 10 subtarefas
- "done" deve ser sempre false

TÍTULO:
%s

DESCRIÇÃO:
%s`, title, description)

	reqBody := map[string]any{
		"model": p.model,
		"input": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "input_text", "text": prompt},
				},
			},
		},
		"text": map[string]any{
			"format": map[string]any{
				"type": "json_object",
			},
		},
	}

	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("openai error: status=%d body=%s", resp.StatusCode, string(rawBytes))
	}

	modelText, err := extractResponsesText(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract model text: %w | raw=%s", err, string(rawBytes))
	}
	if modelText == "" {
		return nil, fmt.Errorf("openai returned empty text | raw=%s", string(rawBytes))
	}

	var out struct {
		Subtasks []SubtaskDTO `json:"subtasks"`
	}
	if err := json.Unmarshal([]byte(modelText), &out); err != nil {
		return nil, fmt.Errorf("failed to parse model json: %w | model_text=%s", err, modelText)
	}

	return out.Subtasks, nil
}
