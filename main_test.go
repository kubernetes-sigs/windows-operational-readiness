package main

import (
	"testing"
)

func TestConfigurationValidation(t *testing.T) {
	var tests = []struct{
		answerWant bool
		categories []string
		category string
	}{
		{
			answerWant: false,
			categories: []string{"category1", "category2"},
			category: "category3",
		},
		{
			answerWant: true,
			categories: []string{"category1", "category2"},
			category: "category2",
		},
	}

	for _, tt := range tests {
		// save global category field
		categoryFlags = tt.categories
		answer := categoryEnabled(tt.category)
		
		if answer != tt.answerWant {
			t.Errorf("got %t, want %t", answer, tt.answerWant)
		}
	}
}
