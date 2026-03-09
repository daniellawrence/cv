package main

import (
	"testing"

	experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"
)

func TestExperienceFields(t *testing.T) {

	tests := []struct {
		id          string
		institution string
		degree      string
	}{
		{"1", "University of Melbourne", "Computer Science"},
		{"2", "RMIT University", "Software Engineering"},
	}

	for _, tt := range tests {

		e := &experiencev1.Experience{
			Id:      tt.id,
			Company: tt.institution,
			Title:   tt.degree,
		}

		if e.Id != tt.id {
			t.Fatalf("expected id %s got %s", tt.id, e.Id)
		}
	}
}
