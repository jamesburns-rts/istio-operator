/*
Copyright 2019 Banzai Cloud.

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

package citadel

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/banzaicloud/istio-operator/pkg/resources/templates"
	"github.com/banzaicloud/istio-operator/pkg/util"
)

func (r *Reconciler) deployment() runtime.Object {
	var args []string
	if util.PointerToBool(r.Config.Spec.SDS.Enabled) {
		args = append(args, "--sds-enabled=true")
	}

	args = append(args,
		"--append-dns-names=true",
		"--grpc-port=8060",
		fmt.Sprintf("--citadel-storage-namespace=%s", r.Config.Namespace),
		fmt.Sprintf("--custom-dns-names=istio-pilot-service-account.%[1]s:istio-pilot.%[1]s", r.Config.Namespace),
		"--monitoring-port=15014",
	)

	if r.Config.Spec.Citadel.CASecretName == "" {
		args = append(args, "--self-signed-ca=true")
	} else {
		args = append(args,
			"--self-signed-ca=false",
			"--signing-cert=/etc/cacerts/ca-cert.pem",
			"--signing-key=/etc/cacerts/ca-key.pem",
			"--root-cert=/etc/cacerts/root-cert.pem",
			"--cert-chain=/etc/cacerts/cert-chain.pem",
		)
	}

	if util.PointerToBool(r.Config.Spec.Citadel.HealthCheck) {
		args = append(args,
			"--liveness-probe-path=/tmp/ca.liveness",
			"--liveness-probe-interval=60s",
			"--probe-check-interval=15s",
		)
	}

	if r.Config.Spec.Citadel.WorkloadCertTTL != "" {
		args = append(args,
			"--workload-cert-ttl",
			r.Config.Spec.Citadel.WorkloadCertTTL,
		)
	}

	if r.Config.Spec.Citadel.MaxWorkloadCertTTL != "" {
		args = append(args,
			"--max-workload-cert-ttl",
			r.Config.Spec.Citadel.MaxWorkloadCertTTL,
		)
	}

	var citadelContainer = apiv1.Container{
		Name:            "citadel",
		Image:           r.Config.Spec.Citadel.Image,
		ImagePullPolicy: r.Config.Spec.ImagePullPolicy,
		Args:            args,
		Resources: templates.GetResourcesRequirementsOrDefault(
			r.Config.Spec.Citadel.Resources,
			r.Config.Spec.DefaultResources,
		),
		TerminationMessagePath:   apiv1.TerminationMessagePathDefault,
		TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
	}

	if util.PointerToBool(r.Config.Spec.Citadel.HealthCheck) {
		citadelContainer.LivenessProbe = &apiv1.Probe{
			Handler: apiv1.Handler{
				Exec: &apiv1.ExecAction{
					Command: []string{
						"/usr/local/bin/istio_ca",
						"probe",
						"--probe-path=/tmp/ca.liveness",
						"--interval=125s",
					},
				},
			},
			InitialDelaySeconds: 60,
			PeriodSeconds:       60,
			FailureThreshold:    30,
			SuccessThreshold:    1,
			TimeoutSeconds:      1,
		}
	}

	if r.Config.Spec.Citadel.CASecretName != "" {
		citadelContainer.VolumeMounts = []apiv1.VolumeMount{
			{
				Name:      "cacerts",
				MountPath: "/etc/cacerts",
				ReadOnly:  true,
			},
		}
	}

	var podSpec = apiv1.PodSpec{
		ServiceAccountName:            serviceAccountName,
		DNSPolicy:                     apiv1.DNSClusterFirst,
		RestartPolicy:                 apiv1.RestartPolicyAlways,
		TerminationGracePeriodSeconds: util.Int64Pointer(int64(30)),
		SecurityContext:               &apiv1.PodSecurityContext{},
		SchedulerName:                 "default-scheduler",
		Containers: []apiv1.Container{
			citadelContainer,
		},
		Affinity:     r.Config.Spec.Citadel.Affinity,
		NodeSelector: r.Config.Spec.Citadel.NodeSelector,
		Tolerations:  r.Config.Spec.Citadel.Tolerations,
	}

	var optional = false
	if r.Config.Spec.Citadel.CASecretName != "" {
		podSpec.Volumes = []apiv1.Volume{
			{
				Name: "cacerts",
				VolumeSource: apiv1.VolumeSource{
					Secret: &apiv1.SecretVolumeSource{
						SecretName:  r.Config.Spec.Citadel.CASecretName,
						Optional:    &optional,
						DefaultMode: util.IntPointer(420),
					},
				},
			},
		}
	}

	var deployment = &appsv1.Deployment{
		ObjectMeta: templates.ObjectMeta(deploymentName, util.MergeLabels(citadelLabels, labelSelector), r.Config),
		Spec: appsv1.DeploymentSpec{
			Replicas: util.IntPointer(1),
			Strategy: templates.DefaultRollingUpdateStrategy(),
			Selector: &metav1.LabelSelector{
				MatchLabels: util.MergeLabels(citadelLabels, labelSelector),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      util.MergeLabels(citadelLabels, labelSelector),
					Annotations: templates.DefaultDeployAnnotations(),
				},
				Spec: podSpec,
			},
		},
	}

	return deployment
}
