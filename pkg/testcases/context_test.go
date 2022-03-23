package testcases

import (
	"github.com/k8sbykeshed/op-readiness/pkg/flags"
	"testing"
)

func TestConfigurationValidation(t *testing.T) {
	var tests = []struct {
		answerWant bool
		categories flags.ArrayFlags
		category   string
	}{
		{
			answerWant: false,
			categories: flags.ArrayFlags{"category1", "category2"},
			category:   "category3",
		},
		{
			answerWant: true,
			categories: flags.ArrayFlags{"category1", "category2"},
			category:   "category2",
		},
	}

	for _, tt := range tests {
		// save global category field
		ctx := TestContext{Categories: tt.categories}
		answer := ctx.CategoryEnabled(tt.category)

		if answer != tt.answerWant {
			t.Errorf("got %t, want %t", answer, tt.answerWant)
		}
	}
}
