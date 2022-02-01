package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type OpTestConfig struct {
	OpTestCases []OpTestCase `yaml:"testCases"`
}

type OpTestCase struct {
	Category           string   `yaml:"category"`
	Focus              []string `yaml:"focus,omitempty"`
	Skip               []string `yaml:"skip,omitempty"`
	KubernetesVersions []string `yaml:"kubernetesVersions,omitempty"`
	WindowsPodImage    string   `yaml:"windows_image,omitempty"`
	LinuxPodImage      string   `yaml:"linux_image,omitempty"`
	Description        string   `yaml:"description,omitempty"`
}

func NewOpTestConfig(inputYamlFile string) (opTestConfig *OpTestConfig) {
	inputFile, err := ioutil.ReadFile(inputYamlFile)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Input yaml file load failed, %v", err))
		return nil
	}
	if err := yaml.Unmarshal(inputFile, &opTestConfig); err != nil {
		zap.L().Error(fmt.Sprintf("Input yaml file unmarshal failed, %v", err))
		return nil
	}
	return
}

func runTest(opTestCase OpTestCase) (string, error) {
	args := []string{
		"--ginkgo.v=true",
		"--ginkgo.debug=true",
		"--kubeconfig=/home/kubo/.kube/config",
		"--ginkgo.focus=%v",
		"--node-os-distro=windows",
		"--ginkgo.skip=%v",
		"--ginkgo.noColor=true",
		"--non-blocking-taints=\"os,node-role.kubernetes.io/master,node.kubernetes.io/not-ready\"",
	}
	argsUsed := fmt.Sprintf(strings.Join(args, " "), opTestCase.Focus, opTestCase.Skip)

	split := strings.Split(argsUsed, " ")

	fmt.Println(argsUsed)
	// TODO(iXinqi): replace the placeholder
	runme := exec.Command("./xxx", split...)
	out, err := runme.CombinedOutput()
	return string(out), err

}

func main() {
	opTestConfig := NewOpTestConfig("./tests.yaml")

	for i, c := range opTestConfig.OpTestCases {
		zap.L().Error(fmt.Sprintf("Starting Operational Readiness Test %v / %v : %v", i, len(opTestConfig.OpTestCases), c.Category))
		// o, e := runTest(c)
		// fmt.Println(o)
		// fmt.Println(e)
		fmt.Println(c.Category, c.Description)
	}
}
