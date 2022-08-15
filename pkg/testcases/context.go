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

	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v3"
	"sigs.k8s.io/windows-operational-readiness/pkg/flags"
)

type TestContext struct {
	E2EBinary  string
	KubeConfig string
	Provider   string
	DryRun     bool
	TestConfig *OpTestConfig
	Categories flags.ArrayFlags
}

func NewTestContext(e2ebinary, kubeconfig, provider string, testConfig *OpTestConfig, dryRun bool, categories flags.ArrayFlags) *TestContext {
	return &TestContext{
		E2EBinary:  e2ebinary,
		KubeConfig: kubeconfig,
		Provider:   provider,
		TestConfig: testConfig,
		DryRun:     dryRun,
		Categories: categories,
	}
}

// CategoryEnabled returns a boolean indicating the test category was passed on flags.
func (o *TestContext) CategoryEnabled(category string) bool {
	// When no categories on flags, set as ALL.
	categories := o.Categories
	if len(categories) == 0 {
		return true
	}

	for _, c := range categories {
		if c == category {
			return true
		}
	}
	return false
}

type OpTestConfig struct {
	OpTestCases []OpTestCase `yaml:"testCases"`
}

// NewOpTestConfig returns the marshalled test config with the tests on it.
func NewOpTestConfig(inputYamlFile string) (opTestConfig *OpTestConfig, err error) {
	inputFile, err := os.ReadFile(inputYamlFile)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Input testcases file load failed, %v", err))
		return nil, err
	}

	if err := yaml.Unmarshal(inputFile, &opTestConfig); err != nil {
		zap.L().Error(fmt.Sprintf("Input testcases file unmarshal failed, %v", err))
		return nil, err
	}

	// Validate YAML configuration.
	if err := opTestConfig.validateYAML(); err != nil {
		return nil, err
	}
	return opTestConfig, nil
}
