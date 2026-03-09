package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	qrcodev1 "github.com/daniellawrence/cv/gen/go/qrcode/v1"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/qrcode", getQRCode)
	return mux
}

func TestGetQRCodeSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/qrcode?url=https://example.com", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	var resp qrcodev1.GenerateQRCodeResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Url != "https://example.com" {
		t.Fatalf("expected url https://example.com got %s", resp.Url)
	}

	if resp.ImageBase64 == "" {
		t.Fatalf("expected base64 image data")
	}

	// verify base64 decodes
	_, err = base64.StdEncoding.DecodeString(resp.ImageBase64)
	if err != nil {
		t.Fatalf("invalid base64 image: %v", err)
	}
}

func TestGetQRCodeMissingURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/qrcode", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 got %d", rr.Code)
	}
}
