package parser

import (
	"encoding/json"
	"fmt"
)

type TokenUsage struct {
	Input  int
	Output int
	Total  int
	Cached int
}

type ClaudeResponse struct {
	Type       string         `json:"type"`
	Role       string         `json:"role"`
	Model      string         `json:"model"`
	Content    []ContentBlock `json:"content"`
	StopReason string         `json:"stop_reason"`
	Usage      ClaudeUsage    `json:"usage"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ClaudeUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

func ParseClaudeJSON(data []byte) (string, TokenUsage, error) {
	var resp ClaudeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", TokenUsage{}, fmt.Errorf("parse claude json: %w", err)
	}

	var output string
	for _, block := range resp.Content {
		if block.Type == "text" {
			if output != "" {
				output += "\n"
			}
			output += block.Text
		}
	}

	tokens := TokenUsage{
		Input:  resp.Usage.InputTokens,
		Output: resp.Usage.OutputTokens,
		Total:  resp.Usage.InputTokens + resp.Usage.OutputTokens,
		Cached: resp.Usage.CacheReadInputTokens,
	}

	return output, tokens, nil
}

type GeminiResponse struct {
	Candidates    []GeminiCandidate `json:"candidates"`
	UsageMetadata GeminiUsage       `json:"usageMetadata"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
	Role  string       `json:"role"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiUsage struct {
	PromptTokenCount        int `json:"promptTokenCount"`
	CandidatesTokenCount    int `json:"candidatesTokenCount"`
	TotalTokenCount         int `json:"totalTokenCount"`
	CachedContentTokenCount int `json:"cachedContentTokenCount"`
}

func ParseGeminiJSON(data []byte) (string, TokenUsage, error) {
	var resp GeminiResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", TokenUsage{}, fmt.Errorf("parse gemini json: %w", err)
	}

	var output string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if output != "" {
				output += "\n"
			}
			output += part.Text
		}
	}

	tokens := TokenUsage{
		Input:  resp.UsageMetadata.PromptTokenCount,
		Output: resp.UsageMetadata.CandidatesTokenCount,
		Total:  resp.UsageMetadata.TotalTokenCount,
		Cached: resp.UsageMetadata.CachedContentTokenCount,
	}

	return output, tokens, nil
}
