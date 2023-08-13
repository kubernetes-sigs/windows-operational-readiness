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
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// validateYAML validate the input YAML and returns the error.
func (s *Specification) validateYAML() error {
	validate := validator.New()
	validate.RegisterStructValidation(SpecificationValidation, Specification{})
	return validate.Struct(s)
}

// SpecificationValidation set the required fields and is used by the validator function.
func SpecificationValidation(sl validator.StructLevel) {
	specification := sl.Current().Interface().(Specification)
	if specification.Category == "" {
		zap.L().Error("Category Required")
		sl.ReportError(specification.Category, "category", "Category", "categoryRequired", "")
	}
	for _, testCase := range specification.TestCases {
		if testCase.Description == "" {
			zap.L().Error("Description Required")
			sl.ReportError(testCase.Description, "description", "Description", "descriptionRequired", "")
		}
		if len(testCase.Focus) == 0 || (len(testCase.Focus) == 1 && len(testCase.Focus[0]) == 0) {
			zap.L().Error("Focus Required")
			sl.ReportError(testCase.Focus, "focus", "Focus", "focusRequired", "")
		}
	}
}
