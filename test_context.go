package main

import (
	"flag"
	"os"

	"k8s.io/client-go/tools/clientcmd"
)

var TestContext TestContextType

type TestContextType struct {
	KubeConfig string
	// Provider identifies the infrastructure provider (gce, gke, aws)
	Provider string
}

func RegisterClusterFlags(flags *flag.FlagSet) {
	flags.StringVar(&TestContext.Provider, "provider", "", "The name of the Kubernetes provider (gce, gke, local, skeleton (the fallback if not set), etc.)")
	flags.StringVar(&TestContext.KubeConfig, clientcmd.RecommendedConfigPathFlag, os.Getenv(clientcmd.RecommendedConfigPathEnvVar), "Path to kubeconfig containing embedded authinfo.")
}
