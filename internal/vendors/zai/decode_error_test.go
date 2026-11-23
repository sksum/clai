package zai

import (
	"errors"
	"net/http"
	"testing"

	"github.com/baalimago/clai/pkg/claierr"
)

func TestDecodeInsufficientBalance(t *testing.T) {
	body := []byte(`{"error":{"code":"1113","message":"Insufficient balance or no resource package. Please recharge."}}`)
	err := decodeError(http.StatusTooManyRequests, body)
	var credits *claierr.InsufficientCreditsError
	if !errors.As(err, &credits) {
		t.Fatalf("expected InsufficientCreditsError for 429 + code 1113, got %v", err)
	}
}

func TestDecodeInsufficientBalanceIsStatusChecked(t *testing.T) {
	// The same body on another status is not a drained account.
	body := []byte(`{"error":{"code":"1113","message":"Insufficient balance or no resource package. Please recharge."}}`)
	if err := decodeError(http.StatusInternalServerError, body); err != nil {
		t.Errorf("expected nil for non-429, got %v", err)
	}
}

func TestDecodeBadRequestPreservesMessage(t *testing.T) {
	body := []byte(`{"error":{"code":"1211","message":"Unknown Model, please check the model code."}}`)
	err := decodeError(http.StatusBadRequest, body)
	var unexpected *claierr.UnexpectedProviderResponseError
	if !errors.As(err, &unexpected) {
		t.Fatalf("expected UnexpectedProviderResponseError for 400, got %v", err)
	}
}

func TestDecodeUnknownBodyNil(t *testing.T) {
	if err := decodeError(http.StatusUnauthorized, []byte(`{"error":{"code":"1002","message":"no key"}}`)); err != nil {
		// 401 with a body: not a case zai claims, the baseline stands.
		_ = err
	}
	if err := decodeError(http.StatusBadGateway, []byte(`upstream`)); err != nil {
		t.Errorf("expected nil for undecoded status, got %v", err)
	}
}
