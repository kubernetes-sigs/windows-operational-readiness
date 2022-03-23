package testcases

import (
	"fmt"
	"github.com/k8sbykeshed/op-readiness/pkg/flags"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type TestContext struct {
	E2EBinary  string
	KubeConfig string
	Provider   string
	TestConfig *OpTestConfig
	Categories flags.ArrayFlags
}

func NewTestContext(e2ebinary, kubeconfig, provider string, testConfig *OpTestConfig, categories flags.ArrayFlags) *TestContext {
	return &TestContext{
		E2EBinary:  e2ebinary,
		KubeConfig: kubeconfig,
		Provider:   provider,
		TestConfig: testConfig,
		Categories: categories,
	}
}

// CategoryEnabled returns a boolean indicating the test category was passed on flags
func (o *TestContext) CategoryEnabled(category string) bool {
	// when no categories on flags, set as ALL
	categories := o.Categories
	if len(categories) <= 0 {
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
	inputFile, err := ioutil.ReadFile(inputYamlFile)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Input testcases file load failed, %v", err))
		return nil, err
	}
	if err := yaml.Unmarshal(inputFile, &opTestConfig); err != nil {
		zap.L().Error(fmt.Sprintf("Input testcases file unmarshal failed, %v", err))
		return nil, err
	}
	// Validate YAML configuration
	if err := opTestConfig.validateYAML(); err != nil {
		return nil, err
	}
	return opTestConfig, nil
}
