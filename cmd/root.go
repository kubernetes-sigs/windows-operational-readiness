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

package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/windows-operational-readiness/pkg/flags"
	"sigs.k8s.io/windows-operational-readiness/pkg/report"
	"sigs.k8s.io/windows-operational-readiness/pkg/testcases"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&testDirectory, "test-directory", "", "Path to YAML root directory containing the tests.")
	rootCmd.PersistentFlags().StringVar(&E2EBinary, "e2e-binary", "./e2e.test", "The E2E Ginkgo default binary used to run the tests.")
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "local", "The name of the Kubernetes provider (gce, gke, aws, local, skeleton, azure etc.)")
	rootCmd.PersistentFlags().StringVar(&kubeConfig, clientcmd.RecommendedConfigPathFlag, os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "Path to kubeconfig containing embedded authinfo.")
	rootCmd.PersistentFlags().StringVar(&reportDir, "report-dir", getEnvOrDefault("ARTIFACTS", ""), "Report dump directory, uses artifact for CI integration when set.")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Do not run actual tests, used for sanity check.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable Ginkgo verbosity.")
	rootCmd.PersistentFlags().Var(&categories, "category", "Append category of tests you want to run, default empty will run all tests.")

	rootCmd.AddCommand(reporterCmd)
}

// NewLoggerConfig return the configuration object for the logger.
func NewLoggerConfig(options ...zap.Option) *zap.Logger {
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		NameKey:     "logger",
		TimeKey:     "timer",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime:  zapcore.RFC3339TimeEncoder,
	}), os.Stdout, zap.InfoLevel)
	return zap.New(core).WithOptions(options...)
}

var (
	E2EBinary     string
	provider      string
	testDirectory string
	kubeConfig    string
	dryRun        bool
	verbose       bool
	reportDir     string
	categories    flags.ArrayFlags

	rootCmd = &cobra.Command{
		Use:   "op-readiness",
		Short: "The Windows Operational Readiness testing suite",
		Long:  "Run this software and make sure your Windows node is suitable for Kubernetes operations.",
		Run: func(cmd *cobra.Command, args []string) {
			zap.ReplaceGlobals(NewLoggerConfig())

			specs, err := testcases.NewOpTestConfig(testDirectory)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Create op-readiness context failed, error is %v", zap.Error(err)))
				os.Exit(1)
			}
			testCtx := testcases.NewTestContext(E2EBinary, kubeConfig, provider, specs, dryRun, verbose, reportDir, categories)
			targetedSpecs := make([]testcases.Specification, 0)
			for _, s := range specs {
				if !testCtx.CategoryEnabled(s.Category) {
					zap.L().Info(fmt.Sprintf("[OpReadinessTests] Skipping Ops Readiness Tests for %s because specification was not specified in category filter", s.Category))
				} else {
					targetedSpecs = append(targetedSpecs, s)
				}
			}

			for specIdx, s := range targetedSpecs {
				zap.L().Info(fmt.Sprintf("[OpReadinessTests] %d / %d Specifications - Running %d Test(s) for Ops Readiness specification: %v", specIdx+1, len(targetedSpecs), len(s.TestCases), s.Category))

				if len(s.TestCases) == 0 {
					zap.L().Info(fmt.Sprintf("[%s] No Operational Readiness tests to run", string(s.Category)))
					continue
				}

				for testIdx, t := range s.TestCases {
					prefix := fmt.Sprintf("%d%d", specIdx+1, testIdx+1)
					logPrefix := fmt.Sprintf("[%s] %v / %v Tests - ", s.Category, testIdx+1, len(s.TestCases))
					zap.L().Info(fmt.Sprintf("%sRunning Operational Readiness Test: %v", logPrefix, t.Description))

					var skipProvider = false
					if len(t.SkipProviders) > 0 {
						for _, p := range t.SkipProviders {
							if p == testCtx.Provider {
								skipProvider = true
							}
						}
					}

					if skipProvider {
						zap.L().Info(fmt.Sprintf("%sSkipping Operational Readiness Test for Provider %v: %v", logPrefix, provider, t.Description))
					} else if err = t.RunTest(testCtx, prefix); err != nil {
						zap.L().Error(fmt.Sprintf("%sFailed Operational Readiness Test: %v, error is %v", logPrefix, t.Description, zap.Error(err)))
					} else {
						zap.L().Info(fmt.Sprintf("%sPassed Operational Readiness Test: %v", logPrefix, t.Description))
					}

					// Specs already ran, cleaning up the XML Junit report
					if testCtx.ReportDir != "" && !testCtx.DryRun {
						fileName := "junit_" + prefix + "01.xml"
						junitReport := path.Join(testCtx.ReportDir, fileName)
						zap.L().Info("Cleaning XML files", zap.String("file", fileName), zap.String("path", junitReport))
						if err := report.CleanupJUnitXML(junitReport, string(s.Category), t.Description, testIdx+1); err != nil {
							zap.L().Error(err.Error())
						}
					}
				}
			}
			zap.L().Info("[OpReadinessTests] Completed running Ops Readiness tests")
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func getEnvOrDefault(key, defaultString string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultString
}
