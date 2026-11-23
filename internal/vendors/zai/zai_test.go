package zai

import (
	"testing"

	"github.com/baalimago/clai/internal/audio/generic"
)

func TestSetupConfigMapping(t *testing.T) {
	v := Default
	tmp := 0.42
	v.FrequencyPenalty = tmp
	mt := 1234
	v.MaxTokens = &mt
	v.Temperature = 0.77
	v.TopP = 0.88
	v.Model = "glm-test-model"

	t.Setenv("ZAI_API_KEY", "any-key")
	if err := v.Setup(); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if v.StreamCompleter.Model != v.Model {
		t.Errorf("expected Model %q, got %q", v.Model, v.StreamCompleter.Model)
	}
	if v.StreamCompleter.FrequencyPenalty == nil || *v.StreamCompleter.FrequencyPenalty != v.FrequencyPenalty {
		t.Errorf("frequency penalty not mapped, got %#v want %v", v.StreamCompleter.FrequencyPenalty, v.FrequencyPenalty)
	}
	if v.StreamCompleter.MaxTokens == nil || *v.StreamCompleter.MaxTokens != *v.MaxTokens {
		t.Errorf("max tokens not mapped, got %#v want %v", v.StreamCompleter.MaxTokens, *v.MaxTokens)
	}
	if v.StreamCompleter.Temperature == nil || *v.StreamCompleter.Temperature != v.Temperature {
		t.Errorf("temperature not mapped, got %#v want %v", v.StreamCompleter.Temperature, v.Temperature)
	}
	if v.StreamCompleter.TopP == nil || *v.StreamCompleter.TopP != v.TopP {
		t.Errorf("top_p not mapped, got %#v want %v", v.StreamCompleter.TopP, v.TopP)
	}
	if v.ToolChoice == nil || *v.ToolChoice != "auto" {
		t.Errorf("tool choice expected 'auto', got %#v", v.ToolChoice)
	}
	if v.StreamCompleter.Vendor != "zai" {
		t.Errorf("expected vendor 'zai', got %q", v.StreamCompleter.Vendor)
	}
	if v.StreamCompleter.DecodeError == nil {
		t.Error("expected DecodeError to be set")
	}
}

func TestTranscriberSetup(t *testing.T) {
	v := TranscriberDefault
	t.Setenv("ZAI_API_KEY", "any-key")
	if err := v.Setup(); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if v.Transcriber.Model != "glm-asr-2512" {
		t.Errorf("expected model glm-asr-2512, got %q", v.Transcriber.Model)
	}
	if v.Transcriber.Format != generic.ResponseFormatJSON {
		t.Errorf("expected plain json format, got %q", v.Transcriber.Format)
	}
	if got := v.Diarization(); got != generic.DiarizationNone {
		t.Errorf("expected no diarization, got %v", got)
	}
}
