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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"sigs.k8s.io/windows-operational-readiness/pkg/flags"
)

type TestContext struct {
	E2EBinary      string
	KubeConfig     string
	Provider       string
	DryRun         bool
	Verbose        bool
	ReportDir      string
	Specifications []Specification
	Categories     flags.ArrayFlags
}

func NewTestContext(e2ebinary, kubeconfig, provider string, specifications []Specification, dryRun bool, verbose bool, reportDir string, categories flags.ArrayFlags) *TestContext {
	return &TestContext{
		E2EBinary:      e2ebinary,
		KubeConfig:     kubeconfig,
		Provider:       provider,
		Specifications: specifications,
		DryRun:         dryRun,
		Verbose:        verbose,
		ReportDir:      reportDir,
		Categories:     categories,
	}
}

// CategoryEnabled returns a boolean indicating the test category was passed on flags.
func (o *TestContext) CategoryEnabled(category Category) bool {
	// When no categories on flags, set as ALL.
	categories := o.Categories
	if len(categories) == 0 {
		return true
	}

	for _, c := range categories {
		if c == string(category) {
			return true
		}
	}
	return false
}

// NewOpTestConfig returns the marshalled test config with the tests on it.
func NewOpTestConfig(testDirectory string) ([]Specification, error) {
	testFiles, err := getTestFiles(testDirectory)
	if err != nil {
		return nil, err
	}

	var categories = make([]string, 0)
	specs := make([]Specification, 0)
	for _, testFile := range testFiles {
		spec, err := NewSpecification(testFile)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Input testcase file load failed for %s. Skipping file. \n%v", testFile, err))
		}
		categories = append(categories, string(spec.Category))
		specs = append(specs, *spec)
	}
	zap.L().Info(fmt.Sprintf("Discovered %d test files for Windows Ops Readiness specifications %s", len(specs), strings.Join(categories, ", ")))
	return specs, nil
}

func getTestFiles(testDirectory string) ([]string, error) {
	if testDirectory == "" {
		ex, err := os.Executable()
		if err != nil {
			return nil, err
		}

		testDirectory, err = filepath.Abs(filepath.Dir(ex))
		if err != nil {
			return nil, err
		}

		testDirectory = testDirectory + "/specifications/"
	}

	testFiles := make([]string, 0)
	dirEntries, err := os.ReadDir(testDirectory)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			testFilePath := testDirectory + dirEntry.Name() + "/spec.yaml"
			testFiles = append(testFiles, testFilePath)
		}
	}
	return testFiles, err
}

func NewSpecification(inputYamlFile string) (spec *Specification, err error) {
	inputFile, err := os.ReadFile(inputYamlFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(inputFile, &spec); err != nil {
		zap.L().Error(fmt.Sprintf("Input testcases file unmarshal failed, %v", err))
		return nil, err
	}

	// Validate YAML configuration.
	if err := spec.validateYAML(); err != nil {
		return nil, err
	}
	return spec, nil
}
