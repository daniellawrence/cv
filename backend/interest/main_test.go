package main

import (
	"testing"

	interestv1 "github.com/daniellawrence/cv/gen/go/interest/v1"
)

func TestInterestFields(t *testing.T) {

	tests := []struct {
		id          string
		institution string
		degree      string
	}{
		{"1", "University of Melbourne", "Computer Science"},
		{"2", "RMIT University", "Software Engineering"},
	}

	for _, tt := range tests {

		e := &interestv1.Interest{
			Id:      tt.id,
			Company: tt.institution,
			Title:   tt.degree,
		}

		if e.Id != tt.id {
			t.Fatalf("expected id %s got %s", tt.id, e.Id)
		}
	}
}
