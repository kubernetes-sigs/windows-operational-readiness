package cmd

import (
	"fmt"
	"github.com/k8sbykeshed/op-readiness/pkg/flags"
	"github.com/k8sbykeshed/op-readiness/pkg/testcases"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

// NewLoggerConfig return the configuration object for the logger
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
	kubeconfig string
	categories flags.ArrayFlags

	rootCmd = &cobra.Command{
		Use:   "op-readiness",
		Short: "The Windows Operational Readiness testing suite",
		Long:  "Run this software and make sure your Windows node is suitable for Kubernetes operations.",
		Run: func(cmd *cobra.Command, args []string) {
			zap.ReplaceGlobals(NewLoggerConfig())

			opTestConfig, err := testcases.NewOpTestConfig(testFile)
			if err != nil {
				log.Fatal(err)
			}
			testCtx := testcases.NewTestContext(E2EBinary, kubeconfig, provider, opTestConfig, categories)

			for i, t := range opTestConfig.OpTestCases {
				if !testCtx.CategoryEnabled(t.Category) {
					continue
				}

				zap.L().Info(fmt.Sprintf("Running Operational Readiness Test %v / %v : %v on %v", i+1, len(opTestConfig.OpTestCases), t.Description, t.Category))
				output, err := t.RunTest(testCtx)
				if err != nil {
					zap.L().Fatal("Test case failed with: ", zap.String("error", err.Error()))
				}
				zap.L().Info(output)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&testFile, "test-file", "tests.yaml", "Path to YAML file containing the tests.")
	rootCmd.PersistentFlags().StringVar(&E2EBinary, "e2e-binary", "./e2e.test", "The E2E Ginkgo default binary used to run the tests.")
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "local", "The name of the Kubernetes provider (gce, gke, local, skeleton (the fallback if not set), etc.)")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, clientcmd.RecommendedConfigPathFlag, os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "Path to kubeconfig containing embedded authinfo.")
	rootCmd.PersistentFlags().Var(&categories, "category", "Append category of tests you want to run, default empty will run all tests.")
}
