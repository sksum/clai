package zai

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/baalimago/clai/pkg/claierr"
)

// Z.ai's error body is {"error":{"code":"...","message":"..."}}; the codes are
// numeric strings. Measured 2026-11-24: 1113 (insufficient balance) arrives as
// HTTP 429, an unknown model or bad parameter as 400, a missing key as 401.
func decodeError(status int, body []byte) error {
	var probe struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	decoded := json.Unmarshal(body, &probe) == nil && strings.TrimSpace(probe.Error.Message) != ""
	switch {
	case status == http.StatusTooManyRequests && probe.Error.Code == "1113":
		return claierr.NewInsufficientCredits(&claierr.APIError{
			StatusCode:   status,
			ProviderCode: probe.Error.Code,
			Body:         string(body),
		})
	case status == http.StatusBadRequest && decoded:
		return claierr.NewUnexpectedProviderResponseWithAPI(&claierr.APIError{
			StatusCode:   status,
			ProviderCode: probe.Error.Code,
			Message:      probe.Error.Message,
			Body:         string(body),
		})
	default:
		return nil
	}
}
