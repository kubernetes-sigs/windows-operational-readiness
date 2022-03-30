package testcases

import (
	"bufio"
	"fmt"
	"os/exec"
)

type OpTestCase struct {
	Category           string   `yaml:"category,omitempty"`
	Focus              []string `yaml:"focus,omitempty"`
	Skip               []string `yaml:"skip,omitempty"`
	KubernetesVersions []string `yaml:"kubernetesVersions,omitempty"`
	WindowsPodImage    string   `yaml:"windows_image,omitempty"`
	LinuxPodImage      string   `yaml:"linux_image,omitempty"`
	Description        string   `yaml:"description,omitempty"`
}

// RunTest runs the binary set in the test context with the parameters from flags
func (o *OpTestCase) RunTest(testCtx *TestContext) {
	cmd := exec.Command(testCtx.E2EBinary,
		"--provider", testCtx.Provider,
		"--kubeconfig", testCtx.KubeConfig,
		"--ginkgo.focus", o.Focus[0],
		"--ginkgo.skip", o.Skip[0],
		"--node-os-distro", "windows",
		"--non-blocking-taints", "os",
	)
	stderr, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}
