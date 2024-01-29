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
	"bufio"
	"io"
	"os/exec"
	"sync"

	"go.uber.org/zap"
)

type Specification struct {
	Category  Category   `yaml:"category,omitempty"`
	TestCases []TestCase `yaml:"testCases,omitempty"`
}

type TestCase struct {
	Description        string   `yaml:"description,omitempty"`
	Focus              []string `yaml:"focus,omitempty"`
	Skip               []string `yaml:"skip,omitempty"`
	KubernetesVersions []string `yaml:"kubernetesVersions,omitempty"` // TODO: If versions are specified, only run tests against specified versions
	SkipProviders      []string `yaml:"skipProviders,omitempty"`      // Test will be skipped for those providers. Use if there is known issues with a particular provider for a given test
}

type Category string

const (
	CoreNetwork           Category = "Core.Network"
	CoreStorage           Category = "Core.Storage"
	CoreScheduling        Category = "Core.Scheduling"
	CoreConcurrent        Category = "Core.Concurrent"
	ExtendHostProcess     Category = "Extend.HostProcess"
	ExtendActiveDirectory Category = "Extend.ActiveDirectory"
	ExtendNetworkPolicy   Category = "Extend.NetworkPolicy"
	ExtendNetwork         Category = "Extend.Network"
	ExtendWorker          Category = "Extend.Worker"
)

// RunTest runs the binary set in the test context with the parameters from flags.
func (t *TestCase) RunTest(testCtx *TestContext, prefix string) error {
	cmd := buildCmd(t, testCtx, prefix)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start and run test command with arguments
	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go redirectOutput(&wg, stdout)
	redirectOutput(nil, stderr)

	wg.Wait()
	return cmd.Wait()
}

func buildCmd(t *TestCase, testCtx *TestContext, prefix string) *exec.Cmd {
	args := []string{
		"--provider", testCtx.Provider,
		"--kubeconfig", testCtx.KubeConfig,
		"--report-dir", testCtx.ReportDir,
		"--report-prefix", prefix,
		"--node-os-distro", "windows",
		"--non-blocking-taints", "os,node-role.kubernetes.io/master,node-role.kubernetes.io/control-plane",
		"--ginkgo.flakeAttempts", "1",
	}

	if testCtx.DryRun {
		args = append(args, "--ginkgo.dryRun")
	}

	if testCtx.Verbose {
		args = append(args, "--ginkgo.v")
		args = append(args, "--ginkgo.trace")
	}

	if len(t.Focus) > 0 {
		focus := t.Focus[0]
		for k, f := range t.Focus {
			if k > 0 {
				focus = focus + "|" + f
			}
		}
		args = append(args, "--ginkgo.focus")
		args = append(args, focus)
	}

	if len(t.Skip) > 0 {
		skip := t.Skip[0]
		for k, s := range t.Skip {
			if k > 0 {
				skip = skip + "|" + s
			}
		}
		args = append(args, "--ginkgo.skip")
		args = append(args, skip)
	}

	return exec.Command(testCtx.E2EBinary, args...)
}

func redirectOutput(wg *sync.WaitGroup, stdout io.ReadCloser) {
	if wg != nil {
		defer wg.Done()
	}
	// Increase max buffer size to 1MB to handle long lines of Ginkgo output and avoid bufio.ErrTooLong errors
	const maxBufferSize = 1024 * 1024
	scanner := bufio.NewScanner(stdout)
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		zap.L().Info(m)
	}
	if err := scanner.Err(); err != nil {
		zap.L().Error(err.Error())
	}
}
