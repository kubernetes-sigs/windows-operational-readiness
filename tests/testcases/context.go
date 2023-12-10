/*
Copyright 2023 The Kubernetes Authors.

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
	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint
	"sigs.k8s.io/windows-operational-readiness/pkg/flags"
	"sigs.k8s.io/windows-operational-readiness/pkg/testcases"
)

var _ = Describe("Test Context", func() {
	Describe("Validating category", func() {
		Context("with specific categories in the flag", func() {
			It("should enable only those categories", func() {
				categories := flags.ArrayFlags{"category1", "category2"}
				testCtx := testcases.TestContext{Categories: categories}

				Expect(testCtx.CategoryEnabled("category3")).To(BeFalse())
				Expect(testCtx.CategoryEnabled("category2")).To(BeTrue())
				Expect(testCtx.CategoryEnabled("category1")).To(BeTrue())
			})
		})
	})
})
