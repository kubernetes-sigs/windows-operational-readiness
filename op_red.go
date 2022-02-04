package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

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

func runTest(opTestCase OpTestCase) (string, error) {
	args := []string{
		"--provider=%v",
		"--kubeconfig=%v",
		"--ginkgo.focus=\"should deny ingress from pods on other namespaces\"",
		"--ginkgo.skip=\"Driver|Slow|Driver\"",
		"--ginkgo.dryRun=true",
		// "--node-os-distro=windows",
	}
	// argsUsed := fmt.Sprintf(strings.Join(args, " "), framework.TestContext.Provider, framework.TestContext.KubeConfig, opTestCase.Focus, opTestCase.Skip)
	argsUsed := fmt.Sprintf(strings.Join(args, " "), TestContext.Provider, TestContext.KubeConfig)

	split := strings.Split(argsUsed, " ")

	fmt.Println(argsUsed)
	runme := exec.Command("./e2e.test", split...)
	out, err := runme.CombinedOutput()
	return string(out), err

}

func handleFlags() {
	// handleFlags sets up all flags and parses the command line.
	RegisterClusterFlags(flag.CommandLine)
	flag.Parse()
}

func main() {

	// Register test flags, then parse flags.
	handleFlags()

	opTestConfig, err := NewOpTestConfig("./tests.yaml")
	if err != nil {
		return
	}

	for i, c := range opTestConfig.OpTestCases {
		zap.L().Error(fmt.Sprintf("Starting Operational Readiness Test %v / %v : %v", i, len(opTestConfig.OpTestCases), c.Category))
		o, e := runTest(c)
		fmt.Println(o)
		fmt.Println(e)
		fmt.Println(c.Category, c.Description)
	}
}
