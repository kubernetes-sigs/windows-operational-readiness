/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
