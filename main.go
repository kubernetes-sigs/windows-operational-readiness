package main

import (
	"github.com/k8sbykeshed/op-readiness/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
