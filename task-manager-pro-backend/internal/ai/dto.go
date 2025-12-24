package ai

import "strings"

type SuggestTitleRequest struct {
	Description string `json:"description"`
}

type SuggestTitleResponse struct {
	Titles []string `json:"titles"`
}

type ImproveDescriptionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ImproveDescriptionResponse struct {
	ImprovedDescription string   `json:"improved_description"`
	Bullets             []string `json:"bullets"`
}

func (r ImproveDescriptionRequest) Validate() string {
	title := strings.TrimSpace(r.Title)
	desc := strings.TrimSpace(r.Description)

	if len(title) < 3 {
		return "title must be at least 3 characters"
	}
	if len(title) > 80 {
		return "title must be at most 80 characters"
	}
	if len(desc) < 10 {
		return "description must be at least 10 characters"
	}
	if len(desc) > 1000 {
		return "description must be at most 1000 characters"
	}
	return ""
}

func (r SuggestTitleRequest) Validate() string {
	desc := strings.TrimSpace(r.Description)
	if len(desc) < 10 {
		return "description must be at least 10 characters"
	}
	if len(desc) > 300 {
		return "description must be at most 300 characters"
	}
	return ""
}
