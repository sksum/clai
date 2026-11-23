// Package zai runs GLM models on Z.ai's OpenAI-compatible API
// (https://docs.z.ai/api-reference/llm/chat-completion). Chat models answer at
// /api/paas/v4/chat/completions and speech-to-text at /api/paas/v4/audio/transcriptions,
// both with the ZAI_API_KEY key.
package zai

import (
	"github.com/baalimago/clai/internal/diag"
	"github.com/baalimago/clai/internal/models"
	"github.com/baalimago/clai/internal/text/generic"
	pub_models "github.com/baalimago/clai/pkg/text/models"
)

const ChatURL = "https://api.z.ai/api/paas/v4/chat/completions"

var Default = GLM{
	Model:       "glm-5.3",
	Temperature: 1.0,
	TopP:        1.0,
	URL:         ChatURL,
}

type GLM struct {
	generic.StreamCompleter
	Model            string  `json:"model"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	MaxTokens        *int    `json:"max_tokens"` // Use a pointer to allow null value
	PresencePenalty  float64 `json:"presence_penalty"`
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
	URL              string  `json:"url"`
}

func (g *GLM) Setup() error {
	err := g.StreamCompleter.Setup("ZAI_API_KEY", ChatURL, "DEBUG_ZAI")
	if err != nil {
		return diag.Errorf("failed to setup stream completer: %w", err)
	}
	g.StreamCompleter.DecodeError = decodeError
	g.StreamCompleter.Model = g.Model
	g.StreamCompleter.FrequencyPenalty = &g.FrequencyPenalty
	g.StreamCompleter.MaxTokens = g.MaxTokens
	g.StreamCompleter.Temperature = &g.Temperature
	g.StreamCompleter.TopP = &g.TopP
	toolChoice := "auto"
	g.ToolChoice = &toolChoice
	g.StreamCompleter.Vendor = "zai"
	g.StreamCompleter.Images = g.ImageCapabilities()
	return nil
}

func (g *GLM) RegisterTool(tool pub_models.LLMTool) {
	g.InternalRegisterTool(tool)
}

// ImageCapabilities is the zai row: GLM takes inline images in chat completions
// and answers tool images the chat-completions way.
func (g *GLM) ImageCapabilities() models.ImageCapabilities {
	return models.ImageCapabilities{
		Input:      models.ImageInputInline,
		ToolResult: models.ToolResultFallback,
		MIMEs:      []string{"image/png", "image/jpeg", "image/gif", "image/webp"},
	}
}
