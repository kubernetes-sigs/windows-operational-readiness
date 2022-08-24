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

// Package lifecycle contains the handlers for the lifecycle hooks.
package handlers

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/util"
	"time"

	sc "github.com/vmware-tanzu/sonobuoy/pkg/client"
	sonodynamic "github.com/vmware-tanzu/sonobuoy/pkg/dynamic"
	"github.com/vmware-tanzu/sonobuoy/pkg/plugin/manifest"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	runtimehooksv1 "sigs.k8s.io/cluster-api/exp/runtime/hooks/api/v1alpha1"

	ctrl "sigs.k8s.io/controller-runtime"
)

const image = "gcr.io/k8s-staging-win-op-rdnss/k8s-win-op-rdnss:latest"

var (
	command = []string{"/app/run.sh"}
	args    = []string{"--e2e-binary", "/app/e2e.test", "--category", "Core.Network", "--dry-run"}
)

// Handler is the handler for the lifecycle hooks.
type Handler struct {
	Config  *rest.Config
	Tracker *remote.ClusterCacheTracker
}

// DoAfterControlPlaneInitialized implements AfterControlPlaneInitialized hook.
func (h *Handler) DoAfterControlPlaneInitialized(ctx context.Context, request *runtimehooksv1.AfterControlPlaneInitializedRequest, response *runtimehooksv1.AfterControlPlaneInitializedResponse) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("AfterControlPlaneInitialized is called, trying to start the Sonobuoy plugin in the WL cluster.")

	if !isWorkloadCPReady(request.Cluster.GetConditions()) {
		response.Status = runtimehooksv1.ResponseStatusFailure
		log.Error(nil, "Cluster Workload Control Plane is not ready yet, retrying.")
		return
	}

	config, err := getWorkloadRESTConfig(ctx, request.Cluster.GetObjectMeta(), h.Tracker)
	if err != nil {
		response.Status = runtimehooksv1.ResponseStatusFailure
		log.Error(err, "Cannot access the workload configuration.")
		return
	}

	if err := h.submitSonobuoyRun(config); err != nil {
		response.Status = runtimehooksv1.ResponseStatusFailure
		log.Error(err, "Cannot setup Sonobuoy testing.")
		return
	}

	log.Info("Sonobuoy tests were dispatched. Check if there's a CNI is installed.", "cluster", request.Cluster.GetName())
	response.Status = runtimehooksv1.ResponseStatusSuccess
	return
}

// getWorkloadClientset extract the configuration from tracker object and generates
// a new clientset based on the workload config
func getWorkloadRESTConfig(ctx context.Context, objectMeta metav1.Object, tracker *remote.ClusterCacheTracker) (*rest.Config, error) {
	// returns clientset
	//return kubernetes.NewForConfig(config)
	objectKey := util.ObjectKey(objectMeta)
	return tracker.GetRESTConfig(ctx, objectKey)
}

// isWorkloadCPReady returns true if the ControlPlaneReady conditions is true, and
// false otherwise.
func isWorkloadCPReady(conditions []clusterv1.Condition) bool {
	for _, condition := range conditions {
		if condition.Type == clusterv1.ControlPlaneReadyCondition && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// submitSonobuoyRun starts the Sonobuoy job, should run in the background.
func (h *Handler) submitSonobuoyRun(config *rest.Config) error {
	var (
		err error
		skc *sonodynamic.APIHelper
		sbc *sc.SonobuoyClient
	)

	if skc, err = sonodynamic.NewAPIHelperFromRESTConfig(config); err != nil {
		return err
	}

	if sbc, err = sc.NewSonobuoyClient(config, skc); err != nil {
		return err
	}

	runCfg := &sc.RunConfig{Wait: time.Duration(0)} // no wait, so we dont block the hook
	runCfg.GenConfig = *generateConfig()
	if err := sbc.Run(runCfg); err != nil {
		return err
	}
	return nil
}

// generateConfig returns a marshalled Sonobuoy configuration
func generateConfig() *sc.GenConfig {
	return &sc.GenConfig{
		EnableRBAC: true, // create RBAC resources.
		StaticPlugins: []*manifest.Manifest{
			{
				Spec: manifest.Container{
					Container: corev1.Container{
						Name:    "op-readiness",
						Command: command,
						Image:   image,
						Args:    args,
						VolumeMounts: []corev1.VolumeMount{
							{
								MountPath: "/tmp/sonobuoy/results",
								Name:      "results",
							},
						},
					},
				},
				SonobuoyConfig: manifest.SonobuoyConfig{
					Driver:       "Job",
					PluginName:   "os-readiness",
					ResultFormat: "raw",
				},
			},
		},
	}
}
