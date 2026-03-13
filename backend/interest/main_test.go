package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	interestv1 "github.com/daniellawrence/cv/gen/go/interest/v1"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/interest", listInterest(nil))
	mux.HandleFunc("/interest/{id}", getInterest(nil))
	return mux
}

func TestListInterest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/interest", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	var resp interestv1.ListInterestResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Interest) == 0 {
		t.Fatalf("expected at least one interest")
	}
}

func TestGetInterestSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/interest/tech", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}

	var interest interestv1.Interest
	err := json.Unmarshal(rr.Body.Bytes(), &interest)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if interest.Id != "tech" {
		t.Fatalf("expected id tech got %s", interest.Id)
	}
}

func TestGetInterestNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/interest/unknown", nil)
	rr := httptest.NewRecorder()

	mux := setupMux()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 got %d", rr.Code)
	}
}
