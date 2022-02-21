package main

import (
	"flag"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var categoryFlags arrayFlags

type TestContextType struct {
	OS         string
	KubeConfig string
	// Provider identifies the infrastructure provider (gce, gke, aws)
	Provider string
}

var TestContext TestContextType

func RegisterClusterFlags(flags *flag.FlagSet) {
	flags.StringVar(&TestContext.OS, "os", "linux", "The OS of running the kubernetes e2e conformance test binary.(darwin or linux)")
	flags.StringVar(&TestContext.Provider, "provider", "", "The name of the Kubernetes provider (gce, gke, local, skeleton (the fallback if not set), etc.)")
	flags.StringVar(&TestContext.KubeConfig, clientcmd.RecommendedConfigPathFlag, os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "Path to kubeconfig containing embedded authinfo.")

	// Tests category flags
	flag.Var(&categoryFlags, "category", "Append category of tests you want to run, default empty will run all tests.")
}
