package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type OpTestConfig struct {
	OpTestCases []OpTestCase `yaml:"testCases"`
}

type OpTestCase struct {
	Category           string   `yaml:"category,omitempty"`
	Focus              []string `yaml:"focus,omitempty"`
	Skip               []string `yaml:"skip,omitempty"`
	KubernetesVersions []string `yaml:"kubernetesVersions,omitempty"`
	WindowsPodImage    string   `yaml:"windows_image,omitempty"`
	LinuxPodImage      string   `yaml:"linux_image,omitempty"`
	Description        string   `yaml:"description,omitempty"`
}

func NewOpTestConfig(inputYamlFile string) (opTestConfig *OpTestConfig, err error) {
	inputFile, err := ioutil.ReadFile(inputYamlFile)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Input yaml file load failed, %v", err))
		return nil, err
	}
	if err := yaml.Unmarshal(inputFile, &opTestConfig); err != nil {
		zap.L().Error(fmt.Sprintf("Input yaml file unmarshal failed, %v", err))
		return nil, err
	}
	validate := validator.New()
	validate.RegisterStructValidation(OpTestConfigValidation, OpTestConfig{})
	if err := validate.Struct(opTestConfig); err != nil {
		return nil, err
	}
	return opTestConfig, nil
}

func OpTestConfigValidation(sl validator.StructLevel) {
	opTestConfig := sl.Current().Interface().(OpTestConfig)

	for _, opTestCase := range opTestConfig.OpTestCases {
		if opTestCase.Category == "" {
			fmt.Println("Category Required")
			sl.ReportError(opTestCase.Category, "category", "Category", "categoryRequired", "")
		}
		if opTestCase.Description == "" {
			fmt.Println("Description Required")
			sl.ReportError(opTestCase.Description, "description", "Description", "descriptionRequired", "")
		}
		if len(opTestCase.Focus) == 0 || (len(opTestCase.Focus) == 1 && len(opTestCase.Focus[0]) == 0) {
			fmt.Println("Focus Required")
			sl.ReportError(opTestCase.Focus, "focus", "Focus", "focusRequired", "")
		}
	}
}

func main() {
	// Register test flags, then parse flags.
	handleFlags()

	opTestConfig, err := NewOpTestConfig("./tests.yaml")
	if err != nil {
		return
	}

	for i, c := range opTestConfig.OpTestCases {
		if !categoryEnabled(c.Category) {
			continue
		}

		zap.L().Error(fmt.Sprintf("Starting Operational Readiness Test %v / %v : %v", i, len(opTestConfig.OpTestCases), c.Category))
		o, e := runTest(c)
		fmt.Println(o)
		fmt.Println(e)
		fmt.Println(c.Category, c.Description)
	}
}

func runTest(opTestCase OpTestCase) (string, error) {
	runme := exec.Command("./e2e_test_binary/"+TestContext.OS+"/e2e.test", "--provider", TestContext.Provider, "--kubeconfig", TestContext.KubeConfig, "--ginkgo.focus", opTestCase.Focus[0], "--ginkgo.skip", opTestCase.Skip[0], "--node-os-distro", "windows", "--non-blocking-taints", "os")
	out, err := runme.CombinedOutput()
	return string(out), err

}

// handleFlags sets up all flag and parses the command line.
func handleFlags() {
	RegisterClusterFlags(flag.CommandLine)
	flag.Parse()
}

// categoryEnabled returns a boolean indicating the test category was passed on flags
func categoryEnabled(category string) bool {
	// when no categories on flags, set as ALL
	if len(categoryFlags) == 0 {
		return true
	}

	for _, cat := range categoryFlags {
		if cat == category {
			return true
		}
	}
	return false
}
