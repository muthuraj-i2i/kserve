/*
Copyright 2021 The KServe Authors.

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

package inferenceservice

import (
	"context"

	"github.com/kserve/kserve/pkg/apis/serving/v1alpha1"
	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/kserve/kserve/pkg/constants"
	"google.golang.org/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TestFixtures contains shared test data and configurations
type TestFixtures struct {
	DefaultResource corev1.ResourceRequirements
	Configs         map[string]string
	Domain          string
}

// NewTestFixtures creates a new instance of shared test fixtures
func NewTestFixtures() *TestFixtures {
	return &TestFixtures{
		DefaultResource: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("1"),
				corev1.ResourceMemory: resource.MustParse("2Gi"),
			},
		},
		Configs: map[string]string{
			"explainers": `{
				"art": {
					"image": "kserve/art-explainer",
					"defaultImageVersion": "latest"
				}
			}`,
			"ingress": `{
				"kserveIngressGateway": "kserve/kserve-ingress-gateway",
				"ingressGateway": "knative-serving/knative-ingress-gateway",
				"localGateway": "knative-serving/knative-local-gateway",
				"localGatewayService": "knative-local-gateway.istio-system.svc.cluster.local"
			}`,
			"storageInitializer": `{
				"image" : "kserve/storage-initializer:latest",
				"memoryRequest": "100Mi",
				"memoryLimit": "1Gi",
				"cpuRequest": "100m",
				"cpuLimit": "1",
				"CaBundleConfigMapName": "",
				"caBundleVolumeMountPath": "/etc/ssl/custom-certs",
				"enableDirectPvcVolumeMount": false
			}`,
		},
		Domain: "example.com",
	}
}

// InferenceServiceBuilder provides a fluent API for building InferenceService objects
type InferenceServiceBuilder struct {
	isvc     *v1beta1.InferenceService
	fixtures *TestFixtures
}

// NewInferenceServiceBuilder creates a new builder for InferenceService objects
func NewInferenceServiceBuilder(name, namespace string, fixtures *TestFixtures) *InferenceServiceBuilder {
	return &InferenceServiceBuilder{
		fixtures: fixtures,
		isvc: &v1beta1.InferenceService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
				Annotations: map[string]string{
					"serving.kserve.io/deploymentMode": "Serverless",
				},
			},
			Spec: v1beta1.InferenceServiceSpec{
				Predictor: v1beta1.PredictorSpec{
					ComponentExtensionSpec: v1beta1.ComponentExtensionSpec{
						MinReplicas: ptr.To(int32(1)),
						MaxReplicas: int32(3),
					},
				},
			},
		},
	}
}

// WithStopAnnotation adds or updates the stop annotation
func (b *InferenceServiceBuilder) WithStopAnnotation(value string) *InferenceServiceBuilder {
	if b.isvc.Annotations == nil {
		b.isvc.Annotations = make(map[string]string)
	}
	b.isvc.Annotations[constants.StopAnnotationKey] = value
	return b
}

// WithInitialScaleAnnotation adds or updates the initial scale annotation
func (b *InferenceServiceBuilder) WithInitialScaleAnnotation(value string) *InferenceServiceBuilder {
	if b.isvc.Annotations == nil {
		b.isvc.Annotations = make(map[string]string)
	}
	b.isvc.Annotations["autoscaling.knative.dev/initial-scale"] = value
	return b
}

// WithAnnotation adds a custom annotation
func (b *InferenceServiceBuilder) WithAnnotation(key, value string) *InferenceServiceBuilder {
	if b.isvc.Annotations == nil {
		b.isvc.Annotations = make(map[string]string)
	}
	b.isvc.Annotations[key] = value
	return b
}

// WithLabel adds a label
func (b *InferenceServiceBuilder) WithLabel(key, value string) *InferenceServiceBuilder {
	if b.isvc.Labels == nil {
		b.isvc.Labels = make(map[string]string)
	}
	b.isvc.Labels[key] = value
	return b
}

// WithTensorflowPredictor adds a TensorFlow predictor
func (b *InferenceServiceBuilder) WithTensorflowPredictor(storageUri, runtimeVersion string) *InferenceServiceBuilder {
	b.isvc.Spec.Predictor.Tensorflow = &v1beta1.TFServingSpec{
		PredictorExtensionSpec: v1beta1.PredictorExtensionSpec{
			StorageURI:     &storageUri,
			RuntimeVersion: proto.String(runtimeVersion),
			Container: corev1.Container{
				Name:      constants.InferenceServiceContainerName,
				Resources: b.fixtures.DefaultResource,
			},
		},
	}
	return b
}

// WithModelPredictor adds a Model predictor
func (b *InferenceServiceBuilder) WithModelPredictor(modelFormat, storageUri, runtimeVersion string) *InferenceServiceBuilder {
	b.isvc.Spec.Predictor.Model = &v1beta1.ModelSpec{
		ModelFormat: v1beta1.ModelFormat{
			Name: modelFormat,
		},
		PredictorExtensionSpec: v1beta1.PredictorExtensionSpec{
			StorageURI:     &storageUri,
			RuntimeVersion: proto.String(runtimeVersion),
		},
	}
	return b
}

// WithTransformer adds a transformer component
func (b *InferenceServiceBuilder) WithTransformer(image string) *InferenceServiceBuilder {
	b.isvc.Spec.Transformer = &v1beta1.TransformerSpec{
		ComponentExtensionSpec: v1beta1.ComponentExtensionSpec{
			MinReplicas: ptr.To(int32(1)),
			MaxReplicas: int32(3),
		},
		PodSpec: v1beta1.PodSpec{
			Containers: []corev1.Container{
				{
					Image:     image,
					Resources: b.fixtures.DefaultResource,
				},
			},
		},
	}
	return b
}

// WithExplainer adds an explainer component
func (b *InferenceServiceBuilder) WithExplainer(explainerType string) *InferenceServiceBuilder {
	storageURI := "s3://test/mnist/explainer"
	b.isvc.Spec.Explainer = &v1beta1.ExplainerSpec{
		ComponentExtensionSpec: v1beta1.ComponentExtensionSpec{
			MinReplicas: ptr.To(int32(1)),
			MaxReplicas: int32(3),
		},
		ART: &v1beta1.ARTExplainerSpec{
			ExplainerExtensionSpec: v1beta1.ExplainerExtensionSpec{
				StorageURI: storageURI,
				Container: corev1.Container{
					Resources: b.fixtures.DefaultResource,
				},
			},
		},
	}
	return b
}

// WithMinReplicas sets the min replicas for the predictor
func (b *InferenceServiceBuilder) WithMinReplicas(minReplicas int32) *InferenceServiceBuilder {
	b.isvc.Spec.Predictor.MinReplicas = &minReplicas
	return b
}

// WithMaxReplicas sets the max replicas for the predictor
func (b *InferenceServiceBuilder) WithMaxReplicas(maxReplicas int32) *InferenceServiceBuilder {
	b.isvc.Spec.Predictor.MaxReplicas = maxReplicas
	return b
}

// Build returns the constructed InferenceService
func (b *InferenceServiceBuilder) Build() *v1beta1.InferenceService {
	return b.isvc.DeepCopy()
}

// ServingRuntimeBuilder provides a fluent API for building ServingRuntime objects
type ServingRuntimeBuilder struct {
	sr       *v1alpha1.ServingRuntime
	fixtures *TestFixtures
}

// NewServingRuntimeBuilder creates a new builder for ServingRuntime objects
func NewServingRuntimeBuilder(name, namespace string, fixtures *TestFixtures) *ServingRuntimeBuilder {
	return &ServingRuntimeBuilder{
		fixtures: fixtures,
		sr: &v1alpha1.ServingRuntime{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: v1alpha1.ServingRuntimeSpec{
				ServingRuntimePodSpec: v1alpha1.ServingRuntimePodSpec{
					Containers: []corev1.Container{
						{
							Name:      constants.InferenceServiceContainerName,
							Resources: fixtures.DefaultResource,
						},
					},
				},
				Disabled: proto.Bool(false),
			},
		},
	}
}

// WithTensorflowSupport adds TensorFlow model format support
func (b *ServingRuntimeBuilder) WithTensorflowSupport(version string) *ServingRuntimeBuilder {
	b.sr.Spec.SupportedModelFormats = append(b.sr.Spec.SupportedModelFormats, v1alpha1.SupportedModelFormat{
		Name:       "tensorflow",
		Version:    proto.String(version),
		AutoSelect: proto.Bool(true),
	})

	// Set default TensorFlow container
	b.sr.Spec.ServingRuntimePodSpec.Containers[0].Image = "tensorflow/serving:1.14.0"
	b.sr.Spec.ServingRuntimePodSpec.Containers[0].Command = []string{"/usr/bin/tensorflow_model_server"}
	b.sr.Spec.ServingRuntimePodSpec.Containers[0].Args = []string{
		"--port=9000",
		"--rest_api_port=8080",
		"--model_base_path=/mnt/models",
		"--rest_api_timeout_in_ms=60000",
	}

	return b
}

// WithImagePullSecrets adds image pull secrets
func (b *ServingRuntimeBuilder) WithImagePullSecrets(secrets ...string) *ServingRuntimeBuilder {
	for _, secret := range secrets {
		b.sr.Spec.ServingRuntimePodSpec.ImagePullSecrets = append(
			b.sr.Spec.ServingRuntimePodSpec.ImagePullSecrets,
			corev1.LocalObjectReference{Name: secret},
		)
	}
	return b
}

// Build returns the constructed ServingRuntime
func (b *ServingRuntimeBuilder) Build() *v1alpha1.ServingRuntime {
	return b.sr.DeepCopy()
}

// ConfigMapBuilder provides a fluent API for building ConfigMap objects
type ConfigMapBuilder struct {
	cm *corev1.ConfigMap
}

// NewConfigMapBuilder creates a new builder for ConfigMap objects
func NewConfigMapBuilder(name, namespace string) *ConfigMapBuilder {
	return &ConfigMapBuilder{
		cm: &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: make(map[string]string),
		},
	}
}

// WithData adds data to the ConfigMap
func (b *ConfigMapBuilder) WithData(data map[string]string) *ConfigMapBuilder {
	for k, v := range data {
		b.cm.Data[k] = v
	}
	return b
}

// WithKServeConfigs adds the standard KServe configuration
func (b *ConfigMapBuilder) WithKServeConfigs(fixtures *TestFixtures) *ConfigMapBuilder {
	return b.WithData(fixtures.Configs)
}

// Build returns the constructed ConfigMap
func (b *ConfigMapBuilder) Build() *corev1.ConfigMap {
	return b.cm.DeepCopy()
}

// TestEnvironment encapsulates commonly needed test resources
type TestEnvironment struct {
	ConfigMap      *corev1.ConfigMap
	ServingRuntime *v1alpha1.ServingRuntime
	Context        context.Context
	Cancel         context.CancelFunc
}

// SetupTestEnvironment creates a test environment with common resources
func SetupTestEnvironment(k8sClient client.Client, namespace string) (*TestEnvironment, error) {
	ctx, cancel := context.WithCancel(context.Background())
	fixtures := NewTestFixtures()

	// Create ConfigMap
	configMap := NewConfigMapBuilder(constants.InferenceServiceConfigMapName, constants.KServeNamespace).
		WithKServeConfigs(fixtures).
		Build()

	if err := k8sClient.Create(ctx, configMap); err != nil {
		cancel()
		return nil, err
	}

	// Create ServingRuntime
	servingRuntime := NewServingRuntimeBuilder("tf-serving", namespace, fixtures).
		WithTensorflowSupport("1").
		Build()

	if err := k8sClient.Create(ctx, servingRuntime); err != nil {
		cancel()
		return nil, err
	}

	return &TestEnvironment{
		ConfigMap:      configMap,
		ServingRuntime: servingRuntime,
		Context:        ctx,
		Cancel:         cancel,
	}, nil
}

// Cleanup removes all resources created in the test environment
func (env *TestEnvironment) Cleanup(k8sClient client.Client) {
	if env.Cancel != nil {
		env.Cancel()
	}

	if env.ConfigMap != nil {
		k8sClient.Delete(env.Context, env.ConfigMap)
	}

	if env.ServingRuntime != nil {
		k8sClient.Delete(env.Context, env.ServingRuntime)
	}
}
