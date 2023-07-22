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

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"

	"sigs.k8s.io/windows-operational-readiness/pkg/flags"
	"sigs.k8s.io/windows-operational-readiness/pkg/testcases"
)

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
	E2EBinary  string
	provider   string
	testFile   string
	kubeConfig string
	dryRun     bool
	verbose    bool
	reportDir  string
	categories flags.ArrayFlags

	rootCmd = &cobra.Command{
		Use:   "op-readiness",
		Short: "The Windows Operational Readiness testing suite",
		Long:  "Run this software and make sure your Windows node is suitable for Kubernetes operations.",
		Run: func(cmd *cobra.Command, args []string) {
			zap.ReplaceGlobals(NewLoggerConfig())

			opTestConfig, err := testcases.NewOpTestConfig(testFile)
			if err != nil {
				zap.L().Error(fmt.Sprintf("Create op-readiness context failed, error is %v", zap.Error(err)))
				os.Exit(1)
			}
			testCtx := testcases.NewTestContext(E2EBinary, kubeConfig, provider, opTestConfig, dryRun, verbose, reportDir, categories)

			for idx, t := range opTestConfig.OpTestCases {
				if !testCtx.CategoryEnabled(t.Category) {
					zap.L().Warn(fmt.Sprintf("[%s] %v / %v - Skipping Operational Readiness Test: %v", t.Category, idx+1, len(opTestConfig.OpTestCases), t.Description))
					continue
				}

				zap.L().Info(fmt.Sprintf("[%s] %v / %v - Running Operational Readiness Test: %v", t.Category, idx+1, len(opTestConfig.OpTestCases), t.Description))
				if err = t.RunTest(testCtx, idx+1); err != nil {
					zap.L().Error(fmt.Sprintf("Operational Readiness Test %v failed, error is %v", t.Description, zap.Error(err)))
				}
			}
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

func init() {
	rootCmd.PersistentFlags().StringVar(&testFile, "test-file", "testcases.yaml", "Path to YAML file containing the tests.")
	rootCmd.PersistentFlags().StringVar(&E2EBinary, "e2e-binary", "./e2e.test", "The E2E Ginkgo default binary used to run the tests.")
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "local", "The name of the Kubernetes provider (gce, gke, aws, local, skeleton, etc.)")
	rootCmd.PersistentFlags().StringVar(&kubeConfig, clientcmd.RecommendedConfigPathFlag, os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "Path to kubeconfig containing embedded authinfo.")
	rootCmd.PersistentFlags().StringVar(&reportDir, "report-dir", getEnvOrDefault("ARTIFACTS", ""), "Report dump directory, uses artifact for CI integration when set.")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Do not run actual tests, used for sanity check.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable Ginkgo verbosity.")
	rootCmd.PersistentFlags().Var(&categories, "category", "Append category of tests you want to run, default empty will run all tests.")
}
