package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	identityv1 "github.com/daniellawrence/cv/gen/go/identity/v1"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/identity", listidentity(nil))
	mux.HandleFunc("/identity/{id}", getIdentity(nil))
	return mux
}

func TestListIdentity(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/identity", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	var resp identityv1.ListIdentityResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Identity) == 0 {
		t.Fatalf("expected at least one identity")
	}
}

func TestGetIdentitySuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/identity/dsl", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	var identity identityv1.Identity
	err := json.Unmarshal(rr.Body.Bytes(), &identity)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if identity.Id != "dsl" {
		t.Fatalf("expected id dsl got %s", identity.Id)
	}
}

func TestGetIdentityNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/identity/unknown", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 got %d", rr.Code)
	}
}
