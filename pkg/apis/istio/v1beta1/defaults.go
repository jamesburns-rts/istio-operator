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

package v1beta1

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/banzaicloud/istio-operator/pkg/util"
)

const (
	defaultImageHub                  = "docker.io/istio"
	defaultImageVersion              = "1.2.5"
	defaultPilotImage                = defaultImageHub + "/" + "pilot" + ":" + defaultImageVersion
	defaultCitadelImage              = defaultImageHub + "/" + "citadel" + ":" + defaultImageVersion
	defaultGalleyImage               = defaultImageHub + "/" + "galley" + ":" + defaultImageVersion
	defaultMixerImage                = defaultImageHub + "/" + "mixer" + ":" + defaultImageVersion
	defaultSidecarInjectorImage      = defaultImageHub + "/" + "sidecar_injector" + ":" + defaultImageVersion
	defaultNodeAgentImage            = defaultImageHub + "/" + "node-agent-k8s" + ":" + defaultImageVersion
	defaultSDSImage                  = defaultImageHub + "/" + "node-agent-k8s" + ":" + defaultImageVersion
	defaultProxyImage                = defaultImageHub + "/" + "proxyv2" + ":" + defaultImageVersion
	defaultProxyInitImage            = defaultImageHub + "/" + "proxy_init" + ":" + defaultImageVersion
	defaultInitCNIImage              = "gcr.io/istio-release/install-cni:master-latest-daily"
	defaultCoreDNSImage              = "coredns/coredns:1.1.2"
	defaultCoreDNSPluginImage        = defaultImageHub + "/coredns-plugin:0.2-istio-1.1"
	defaultIncludeIPRanges           = "*"
	defaultReplicaCount              = 1
	defaultMinReplicas               = 1
	defaultMaxReplicas               = 5
	defaultTraceSampling             = 1.0
	defaultIngressGatewayServiceType = apiv1.ServiceTypeLoadBalancer
	defaultEgressGatewayServiceType  = apiv1.ServiceTypeClusterIP
	outboundTrafficPolicyAllowAny    = "ALLOW_ANY"
	defaultZipkinAddress             = "zipkin.%s:9411"
	defaultInitCNIBinDir             = "/opt/cni/bin"
	defaultInitCNIConfDir            = "/etc/cni/net.d"
	defaultInitCNILogLevel           = "info"
	defaultImagePullPolicy           = "IfNotPresent"
	defaultMeshExpansion             = false
	ingress                          = "ingress"
	egress                           = "egress"
)

var defaultResources = &apiv1.ResourceRequirements{
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU: resource.MustParse("10m"),
	},
}

var defaultProxyResources = &apiv1.ResourceRequirements{
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("100m"),
		apiv1.ResourceMemory: resource.MustParse("128Mi"),
	},
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("2000m"),
		apiv1.ResourceMemory: resource.MustParse("1024Mi"),
	},
}

var defaultInitResources = &apiv1.ResourceRequirements{
	Requests: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("10m"),
		apiv1.ResourceMemory: resource.MustParse("10Mi"),
	},
	Limits: apiv1.ResourceList{
		apiv1.ResourceCPU:    resource.MustParse("100m"),
		apiv1.ResourceMemory: resource.MustParse("50Mi"),
	},
}

var defaultIngressGatewayPorts = []apiv1.ServicePort{
	{Port: 15020, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15020), Name: "status-port", NodePort: 31460},
	{Port: 80, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(80), Name: "http2", NodePort: 31380},
	{Port: 443, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(443), Name: "https", NodePort: 31390},
	{Port: 15443, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15443), Name: "tls", NodePort: 31450},
}

var defaultEgressGatewayPorts = []apiv1.ServicePort{
	{Port: 80, Name: "http2", Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(80)},
	{Port: 443, Name: "https", Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(443)},
	{Port: 15443, Protocol: apiv1.ProtocolTCP, TargetPort: intstr.FromInt(15443), Name: "tls"},
}

func SetDefaults(config *Istio) {
	if config.Spec.IncludeIPRanges == "" {
		config.Spec.IncludeIPRanges = defaultIncludeIPRanges
	}
	if config.Spec.Proxy.Resources == nil {
		if config.Spec.DefaultResources == nil {
			config.Spec.Proxy.Resources = defaultProxyResources
		} else {
			config.Spec.Proxy.Resources = defaultResources
		}
	}
	if config.Spec.DefaultResources == nil {
		config.Spec.DefaultResources = defaultResources
	}
	// Pilot config
	if config.Spec.Pilot.Enabled == nil {
		config.Spec.Pilot.Enabled = util.BoolPointer(true)
	}
	if config.Spec.Pilot.Image == "" {
		config.Spec.Pilot.Image = defaultPilotImage
	}
	if config.Spec.Pilot.Sidecar == nil {
		config.Spec.Pilot.Sidecar = util.BoolPointer(true)
	}
	if config.Spec.Pilot.ReplicaCount == 0 {
		config.Spec.Pilot.ReplicaCount = defaultReplicaCount
	}
	if config.Spec.Pilot.MinReplicas == 0 {
		config.Spec.Pilot.MinReplicas = defaultMinReplicas
	}
	if config.Spec.Pilot.MaxReplicas == 0 {
		config.Spec.Pilot.MaxReplicas = defaultMaxReplicas
	}
	if config.Spec.Pilot.TraceSampling == 0 {
		config.Spec.Pilot.TraceSampling = defaultTraceSampling
	}
	// Citadel config
	if config.Spec.Citadel.Enabled == nil {
		config.Spec.Citadel.Enabled = util.BoolPointer(true)
	}
	if config.Spec.Citadel.Image == "" {
		config.Spec.Citadel.Image = defaultCitadelImage
	}
	// Galley config
	if config.Spec.Galley.Enabled == nil {
		config.Spec.Galley.Enabled = util.BoolPointer(true)
	}
	if config.Spec.Galley.Image == "" {
		config.Spec.Galley.Image = defaultGalleyImage
	}
	if config.Spec.Galley.ReplicaCount == 0 {
		config.Spec.Galley.ReplicaCount = defaultReplicaCount
	}
	// Gateways config
	if config.Spec.Gateways.Enabled == nil {
		config.Spec.Gateways.Enabled = util.BoolPointer(true)
	}
	if config.Spec.Gateways.Configs == nil {
		config.Spec.Gateways.Configs = make(map[string]*GatewayConfiguration)
	}
	if ingressConf := config.Spec.Gateways.Configs[ingress]; ingressConf == nil {
		config.Spec.Gateways.Configs[ingress] = new(GatewayConfiguration)
	}
	if egressConf := config.Spec.Gateways.Configs[egress]; egressConf == nil {
		config.Spec.Gateways.Configs[egress] = new(GatewayConfiguration)
	}
	for key, conf := range config.Spec.Gateways.Configs {
		if conf.Enabled == nil {
			conf.Enabled = util.BoolPointer(true)
		}
		if conf.ReplicaCount == 0 {
			conf.ReplicaCount = defaultReplicaCount
		}
		if conf.MinReplicas == 0 {
			conf.MinReplicas = defaultMinReplicas
		}
		if conf.MaxReplicas == 0 {
			conf.MaxReplicas = defaultMaxReplicas
		}
		if conf.SDS.Enabled == nil {
			conf.SDS.Enabled = util.BoolPointer(false)
		}
		if conf.SDS.Image == "" {
			conf.SDS.Image = defaultSDSImage
		}
		if conf.ServiceType == "" {
			switch key {
			case egress:
				conf.ServiceType = defaultEgressGatewayServiceType
			default:
				conf.ServiceType = defaultIngressGatewayServiceType
			}
		}
		if len(conf.Ports) == 0 {
			switch key {
			case egress:
				conf.Ports = defaultEgressGatewayPorts
			case ingress:
				conf.Ports = defaultIngressGatewayPorts
			}
		}
	}
	// Mixer config
	if config.Spec.Mixer.Enabled == nil {
		config.Spec.Mixer.Enabled = util.BoolPointer(true)
	}
	if config.Spec.Mixer.Image == "" {
		config.Spec.Mixer.Image = defaultMixerImage
	}
	if config.Spec.Mixer.ReplicaCount == 0 {
		config.Spec.Mixer.ReplicaCount = defaultReplicaCount
	}
	if config.Spec.Mixer.MinReplicas == 0 {
		config.Spec.Mixer.MinReplicas = defaultMinReplicas
	}
	if config.Spec.Mixer.MaxReplicas == 0 {
		config.Spec.Mixer.MaxReplicas = defaultMaxReplicas
	}
	// SidecarInjector config
	if config.Spec.SidecarInjector.Enabled == nil {
		config.Spec.SidecarInjector.Enabled = util.BoolPointer(true)
	}
	if config.Spec.SidecarInjector.AutoInjectionPolicyEnabled == nil {
		config.Spec.SidecarInjector.AutoInjectionPolicyEnabled = util.BoolPointer(true)
	}
	if config.Spec.SidecarInjector.Image == "" {
		config.Spec.SidecarInjector.Image = defaultSidecarInjectorImage
	}
	if config.Spec.SidecarInjector.ReplicaCount == 0 {
		config.Spec.SidecarInjector.ReplicaCount = defaultReplicaCount
	}
	if config.Spec.SidecarInjector.InitCNIConfiguration.Enabled == nil {
		config.Spec.SidecarInjector.InitCNIConfiguration.Enabled = util.BoolPointer(false)
	}
	if config.Spec.SidecarInjector.InitCNIConfiguration.Image == "" {
		config.Spec.SidecarInjector.InitCNIConfiguration.Image = defaultInitCNIImage
	}
	if config.Spec.SidecarInjector.InitCNIConfiguration.BinDir == "" {
		config.Spec.SidecarInjector.InitCNIConfiguration.BinDir = defaultInitCNIBinDir
	}
	if config.Spec.SidecarInjector.InitCNIConfiguration.ConfDir == "" {
		config.Spec.SidecarInjector.InitCNIConfiguration.ConfDir = defaultInitCNIConfDir
	}
	if config.Spec.SidecarInjector.InitCNIConfiguration.ExcludeNamespaces == nil {
		config.Spec.SidecarInjector.InitCNIConfiguration.ExcludeNamespaces = []string{config.Namespace}
	}
	if config.Spec.SidecarInjector.InitCNIConfiguration.LogLevel == "" {
		config.Spec.SidecarInjector.InitCNIConfiguration.LogLevel = defaultInitCNILogLevel
	}
	if config.Spec.SidecarInjector.Init.Resources == nil {
		config.Spec.SidecarInjector.Init.Resources = defaultInitResources
	}
	// SDS config
	if config.Spec.SDS.Enabled == nil {
		config.Spec.SDS.Enabled = util.BoolPointer(false)
	}
	// NodeAgent config
	if config.Spec.NodeAgent.Enabled == nil {
		config.Spec.NodeAgent.Enabled = util.BoolPointer(false)
	}
	if config.Spec.NodeAgent.Image == "" {
		config.Spec.NodeAgent.Image = defaultNodeAgentImage
	}
	// Proxy config
	if config.Spec.Proxy.Image == "" {
		config.Spec.Proxy.Image = defaultProxyImage
	}
	// Proxy Init config
	if config.Spec.ProxyInit.Image == "" {
		config.Spec.ProxyInit.Image = defaultProxyInitImage
	}
	if config.Spec.Proxy.ComponentLogLevel == "" {
		config.Spec.Proxy.ComponentLogLevel = "misc:error"
	}
	if config.Spec.Proxy.LogLevel == "" {
		config.Spec.Proxy.LogLevel = "warning"
	}
	if config.Spec.Proxy.DNSRefreshRate == "" {
		config.Spec.Proxy.DNSRefreshRate = "300s"
	}
	// PDB config
	if config.Spec.DefaultPodDisruptionBudget.Enabled == nil {
		config.Spec.DefaultPodDisruptionBudget.Enabled = util.BoolPointer(false)
	}
	// Outbound traffic policy config
	if config.Spec.OutboundTrafficPolicy.Mode == "" {
		config.Spec.OutboundTrafficPolicy.Mode = outboundTrafficPolicyAllowAny
	}
	// Tracing config
	if config.Spec.Tracing.Enabled == nil {
		config.Spec.Tracing.Enabled = util.BoolPointer(true)
	}
	if config.Spec.Tracing.Tracer == "" {
		config.Spec.Tracing.Tracer = TracerTypeZipkin
	}
	if config.Spec.Tracing.Datadog.Address == "" {
		config.Spec.Tracing.Datadog.Address = "$(HOST_IP):8126"
	}
	if config.Spec.Tracing.Zipkin.Address == "" {
		config.Spec.Tracing.Zipkin.Address = fmt.Sprintf(defaultZipkinAddress, config.Namespace)
	}

	// Multi mesh support
	if config.Spec.MultiMesh == nil {
		config.Spec.MultiMesh = util.BoolPointer(false)
	}

	// Istio CoreDNS for multi mesh support
	if config.Spec.IstioCoreDNS.Enabled == nil {
		config.Spec.IstioCoreDNS.Enabled = util.BoolPointer(false)
	}
	if config.Spec.IstioCoreDNS.Image == "" {
		config.Spec.IstioCoreDNS.Image = defaultCoreDNSImage
	}
	if config.Spec.IstioCoreDNS.PluginImage == "" {
		config.Spec.IstioCoreDNS.PluginImage = defaultCoreDNSPluginImage
	}
	if config.Spec.IstioCoreDNS.ReplicaCount == 0 {
		config.Spec.IstioCoreDNS.ReplicaCount = defaultReplicaCount
	}

	if config.Spec.ImagePullPolicy == "" {
		config.Spec.ImagePullPolicy = defaultImagePullPolicy
	}

	if config.Spec.MeshExpansion == nil {
		config.Spec.MeshExpansion = util.BoolPointer(defaultMeshExpansion)
	}

	if config.Spec.UseMCP == nil {
		config.Spec.UseMCP = util.BoolPointer(true)
	}
}

func SetRemoteIstioDefaults(remoteconfig *RemoteIstio) {
	if remoteconfig.Spec.IncludeIPRanges == "" {
		remoteconfig.Spec.IncludeIPRanges = defaultIncludeIPRanges
	}
	// SidecarInjector config
	if remoteconfig.Spec.SidecarInjector.ReplicaCount == 0 {
		remoteconfig.Spec.SidecarInjector.ReplicaCount = defaultReplicaCount
	}
}
