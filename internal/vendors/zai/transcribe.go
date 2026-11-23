package zai

import (
	"time"

	"github.com/baalimago/clai/internal/audio/generic"
	"github.com/baalimago/clai/internal/diag"
)

const TranscribeURL = "https://api.z.ai/api/paas/v4/audio/transcriptions"

// MaxRequest is the longest clip one glm-asr request may carry: the API takes
// at most 30 seconds and 25 MB of wav or mp3 (docs.z.ai, Audio Transcriptions).
// A second under keeps clear of the boundary the splitter cuts at.
const MaxRequest = 29 * time.Second

var TranscriberDefault = Transcriber{
	Model: "glm-asr-2512",
}

type Transcriber struct {
	generic.Transcriber
	Model string `json:"model"`
}

func (t *Transcriber) Setup() error {
	err := t.Transcriber.Setup("ZAI_API_KEY", TranscribeURL)
	if err != nil {
		return diag.Errorf("failed to setup transcriber: %w", err)
	}
	t.Transcriber.Model = t.Model
	t.Transcriber.DecodeError = decodeError
	// glm-asr answers OpenAI's plain json: the whole text in one field, no
	// segments, so a request is parsed as a single segment.
	t.Transcriber.Format = generic.ResponseFormatJSON
	return nil
}

// Diarization is none: glm-asr returns plain text, so speakers are named by the
// agent from context and the voice library instead of the model.
func (t *Transcriber) Diarization() generic.Diarization {
	return generic.DiarizationNone
}
