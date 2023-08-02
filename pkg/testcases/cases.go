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
	"encoding/xml"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"

	"go.uber.org/zap"

	"sigs.k8s.io/windows-operational-readiness/pkg/report"
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
func (t *TestCase) RunTest(testCtx *TestContext, idx int) error {
	cmd := buildCmd(t, testCtx, idx)

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

	redirectOutput(stdout)
	redirectOutput(stderr)

	if err := cmd.Wait(); err != nil {
		return err
	}

	if testCtx.ReportDir != "" && !testCtx.DryRun {
		fileName := "junit_" + strconv.Itoa(idx) + "01.xml"
		junitReport := path.Join(testCtx.ReportDir, fileName)
		zap.L().Info("Cleaning XML files", zap.String("file", fileName), zap.String("path", junitReport))
		if err := CleanupJUnitXML(junitReport); err != nil {
			return err
		}
	}
	return nil
}

func buildCmd(t *TestCase, testCtx *TestContext, idx int) *exec.Cmd {
	args := []string{
		"--provider", testCtx.Provider,
		"--kubeconfig", testCtx.KubeConfig,
		"--report-dir", testCtx.ReportDir,
		"--report-prefix", strconv.Itoa(idx),
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

// CleanupJUnitXML removes unnecessary skipped tests from the report.
func CleanupJUnitXML(path string) error {
	zap.L().Info("Getting XML content from file", zap.String("path", path))
	content, err := getXMLContent(path)
	if err != nil {
		return err
	}

	var testSuites report.TestSuites
	if err = xml.Unmarshal(content, &testSuites); err != nil {
		return err
	}
	for i, suite := range testSuites.Suites {
		var cleanTests []report.TestCase
		for _, test := range suite.TestCases {
			if test.Status != report.StatusSkipped &&
				test.Name != "[SynchronizedBeforeSuite]" &&
				test.Name != "[SynchronizedAfterSuite]" &&
				test.Name != "[ReportAfterSuite] Kubernetes e2e suite report" {
				cleanTests = append(cleanTests, test)
			}
		}
		testSuites.Suites[i].TestCases = cleanTests
		zap.L().Info("Saving cleaned tests.", zap.Int("number", len(cleanTests)))
	}
	// write back the cleaned up YAML to a writer
	var cleanContent []byte
	if cleanContent, err = xml.MarshalIndent(testSuites, "  ", "    "); err != nil {
		return err
	}
	if err = writeFileContent(path, cleanContent); err != nil {
		return err
	}
	return nil
}

// writeFileContent save the content to a file.
func writeFileContent(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

// getXMLContent returns the content in bytes of an existent file.
func getXMLContent(path string) ([]byte, error) {
	var err error
	// check if file exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return []byte{}, err
	}

	var file *os.File
	if file, err = os.Open(path); err != nil {
		return []byte{}, err
	}
	defer file.Close()

	var content []byte
	if content, err = io.ReadAll(file); err != nil {
		return []byte{}, err
	}

	return content, nil
}

func redirectOutput(stdout io.ReadCloser) {
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		zap.L().Info(m)
	}
}
