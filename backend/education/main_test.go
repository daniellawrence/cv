package main

import (
	"testing"

	educationv1 "github.com/daniellawrence/cv/gen/go/education/v1"
)

func TestEducationFields(t *testing.T) {

	tests := []struct {
		id          string
		institution string
		degree      string
	}{
		{"1", "University of Melbourne", "Computer Science"},
		{"2", "RMIT University", "Software Engineering"},
	}

	for _, tt := range tests {

		e := &educationv1.Education{
			Id:          tt.id,
			Institution: tt.institution,
			Degree:      tt.degree,
		}

		if e.Id != tt.id {
			t.Fatalf("expected id %s got %s", tt.id, e.Id)
		}
	}
}
