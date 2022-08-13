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

package main

import (
	"flag"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/windows-operational-readiness/exp/handlers"

	runtimecatalog "sigs.k8s.io/cluster-api/exp/runtime/catalog"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"
	"sigs.k8s.io/cluster-api/exp/runtime/server"
)

var (
	catalog  = runtimecatalog.New()
	setupLog = ctrl.Log.WithName("setup")

	// Flags.
	profilerAddress string
	webhookPort     int
	webhookCertDir  string
	logOptions      = logs.NewOptions()
)

func init() {
	// Register the Runtime Hook types into the catalog.
	_ = runtimehooksv1.AddToCatalog(catalog)
}

// InitFlags initializes the flags.
func InitFlags(fs *pflag.FlagSet) {
	logs.AddFlags(fs, logs.SkipLoggingConfigurationFlags())
	logOptions.AddFlags(fs)

	fs.StringVar(&profilerAddress, "profiler-address", "",
		"Bind address to expose the pprof profiler (e.g. localhost:6060)")

	fs.IntVar(&webhookPort, "webhook-port", 9443,
		"Webhook Server port")

	fs.StringVar(&webhookCertDir, "webhook-cert-dir", "/tmp/k8s-webhook-server/serving-certs/",
		"Webhook cert dir, only used when webhook-port is specified.")
}

func main() {
	InitFlags(pflag.CommandLine)
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := logOptions.ValidateAndApply(nil); err != nil {
		setupLog.Error(err, "unable to start extension")
		os.Exit(1)
	}

	// klog.Background will automatically use the right logger.
	ctrl.SetLogger(klog.Background())

	if profilerAddress != "" {
		klog.Infof("Profiler listening for requests at %s", profilerAddress)
		go func() {
			klog.Info(http.ListenAndServe(profilerAddress, nil))
		}()
	}

	ctx := ctrl.SetupSignalHandler()

	webhookServer, err := server.NewServer(server.Options{
		Catalog: catalog,
		Port:    webhookPort,
		CertDir: webhookCertDir,
	})
	if err != nil {
		setupLog.Error(err, "error creating webhook server")
		os.Exit(1)
	}

	registerHooks(webhookServer)

	setupLog.Info("Starting Runtime Extension server")
	if err := webhookServer.Start(ctx); err != nil {
		setupLog.Error(err, "error running webhook server")
		os.Exit(1)
	}
}

// registerHooks register extension handlers.
func registerHooks(webhookServer *server.Server) {
	restConfig, err := ctrl.GetConfig()
	if err != nil {
		setupLog.Error(err, "error getting config for the cluster")
		os.Exit(1)
	}

	c, err := client.New(restConfig, client.Options{})
	if err != nil {
		setupLog.Error(err, "error creating client to the cluster")
		os.Exit(1)
	}

	handler := handlers.Handler{Client: c}

	if err := webhookServer.AddExtensionHandler(server.ExtensionHandler{
		Hook:           runtimehooksv1.AfterControlPlaneInitialized,
		Name:           "after-controlplane-initialized",
		HandlerFunc:    handler.DoAfterControlPlaneInitialized,
		TimeoutSeconds: pointer.Int32(5),
		FailurePolicy:  toPtr(runtimehooksv1.FailurePolicyFail),
	}); err != nil {
		setupLog.Error(err, "error adding handler")
		os.Exit(1)
	}
}

func toPtr(f runtimehooksv1.FailurePolicy) *runtimehooksv1.FailurePolicy {
	return &f
}
