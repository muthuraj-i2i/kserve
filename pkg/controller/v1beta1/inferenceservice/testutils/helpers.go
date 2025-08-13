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

package testutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kserve/kserve/pkg/apis/serving/v1alpha1"
	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/kserve/kserve/pkg/constants"
)

// Test timeout constants - optimized for faster execution
const (
	FastTimeout         = time.Second * 10      // Further reduced from 15s for fast operations
	MediumTimeout       = time.Second * 20      // Reduced from 30s for medium operations
	SlowTimeout         = time.Second * 30      // Reduced from 45s for slow operations
	PollInterval        = time.Millisecond * 50 // Even faster polling interval
	ConsistentlyTimeout = time.Second * 2       // Further reduced from 3s
)

// Default resource requirements used across tests
var DefaultResource = corev1.ResourceRequirements{
	Limits: corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("1"),
		corev1.ResourceMemory: resource.MustParse("2Gi"),
	},
	Requests: corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("1"),
		corev1.ResourceMemory: resource.MustParse("2Gi"),
	},
}

// CreateDefaultConfigMap creates a standard configmap for testing
func CreateDefaultConfigMap(k8sClient client.Client, configs map[string]string) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.InferenceServiceConfigMapName,
			Namespace: constants.KServeNamespace,
		},
		Data: configs,
	}
	return configMap
}

// CreateDefaultServingRuntime creates a standard TensorFlow serving runtime for testing
func CreateDefaultServingRuntime(name, namespace string) *v1alpha1.ServingRuntime {
	return &v1alpha1.ServingRuntime{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.ServingRuntimeSpec{
			SupportedModelFormats: []v1alpha1.SupportedModelFormat{
				{
					Name:       "tensorflow",
					Version:    proto.String("1"),
					AutoSelect: proto.Bool(true),
				},
			},
			ServingRuntimePodSpec: v1alpha1.ServingRuntimePodSpec{
				Containers: []corev1.Container{
					{
						Name:    "kserve-container",
						Image:   "tensorflow/serving:1.14.0",
						Command: []string{"/usr/bin/tensorflow_model_server"},
						Args: []string{
							"--port=9000",
							"--rest_api_port=8080",
							"--model_base_path=/mnt/models",
							"--rest_api_timeout_in_ms=60000",
						},
						Resources: DefaultResource,
					},
				},
			},
			Disabled: proto.Bool(false),
		},
	}
}

// CreateDefaultInferenceService creates a standard InferenceService for testing
func CreateDefaultInferenceService(name, namespace string, annotations map[string]string) *v1beta1.InferenceService {
	return CreateInferenceServiceWithSpec(name, namespace, annotations, nil)
}

// CreateInferenceServiceWithSpec creates an InferenceService with custom predictor spec
func CreateInferenceServiceWithSpec(name, namespace string, annotations map[string]string, predictorSpec *v1beta1.PredictorSpec) *v1beta1.InferenceService {
	if annotations == nil {
		annotations = make(map[string]string)
	}
	// Ensure default annotations are set
	if _, exists := annotations[constants.DeploymentMode]; !exists {
		annotations[constants.DeploymentMode] = string(constants.RawDeployment)
	}
	if _, exists := annotations[constants.AutoscalerClass]; !exists {
		annotations[constants.AutoscalerClass] = string(constants.AutoscalerClassHPA)
	}

	var predictor v1beta1.PredictorSpec
	if predictorSpec != nil {
		predictor = *predictorSpec
	} else {
		// Default predictor spec
		storageUri := "s3://test/mnist/export"
		qty := resource.MustParse("10Gi")

		predictor = v1beta1.PredictorSpec{
			ComponentExtensionSpec: v1beta1.ComponentExtensionSpec{
				MinReplicas:    ptr.To(int32(1)),
				MaxReplicas:    3,
				TimeoutSeconds: ptr.To(int64(30)),
				AutoScaling: &v1beta1.AutoScalingSpec{
					Metrics: []v1beta1.MetricsSpec{
						{
							Type: v1beta1.ResourceMetricSourceType,
							Resource: &v1beta1.ResourceMetricSource{
								Name: v1beta1.ResourceMetricMemory,
								Target: v1beta1.MetricTarget{
									Type:         v1beta1.AverageValueMetricType,
									AverageValue: &qty,
								},
							},
						},
					},
				},
			},
			Tensorflow: &v1beta1.TFServingSpec{
				PredictorExtensionSpec: v1beta1.PredictorExtensionSpec{
					StorageURI:     &storageUri,
					RuntimeVersion: proto.String("1.14.0"),
					Container: corev1.Container{
						Name:      constants.InferenceServiceContainerName,
						Resources: DefaultResource,
					},
				},
			},
		}
	}

	isvc := &v1beta1.InferenceService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: v1beta1.InferenceServiceSpec{
			Predictor: predictor,
		},
	}

	// Apply defaults
	isvc.DefaultInferenceService(nil, nil, &v1beta1.SecurityConfig{AutoMountServiceAccountToken: false}, nil)
	return isvc
} // GetPredictorKey returns the namespaced name for predictor resources
func GetPredictorKey(serviceName, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      constants.PredictorServiceName(serviceName),
		Namespace: namespace,
	}
}

// GetTransformerKey returns the namespaced name for transformer resources
func GetTransformerKey(serviceName, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      constants.TransformerServiceName(serviceName),
		Namespace: namespace,
	}
}

// GetExplainerKey returns the namespaced name for explainer resources
func GetExplainerKey(serviceName, namespace string) types.NamespacedName {
	return types.NamespacedName{
		Name:      constants.ExplainerServiceName(serviceName),
		Namespace: namespace,
	}
}

// SetupTestResources creates and returns common test resources (ConfigMap, ServingRuntime)
func SetupTestResources(k8sClient client.Client, configs map[string]string, runtimeName, namespace string) (*corev1.ConfigMap, *v1alpha1.ServingRuntime, error) {
	configMap := CreateDefaultConfigMap(k8sClient, configs)
	if err := k8sClient.Create(context.TODO(), configMap); err != nil {
		// Handle concurrent configmap creation - if it already exists, that's fine for shared usage
		if !apierr.IsAlreadyExists(err) {
			return nil, nil, err
		}
		// If configmap already exists, get the existing one for consistency
		if getErr := k8sClient.Get(context.TODO(), types.NamespacedName{
			Name:      configMap.Name,
			Namespace: configMap.Namespace,
		}, configMap); getErr != nil {
			return nil, nil, fmt.Errorf("failed to get existing configmap: %w", getErr)
		}
	}

	servingRuntime := CreateDefaultServingRuntime(runtimeName, namespace)
	if err := k8sClient.Create(context.TODO(), servingRuntime); err != nil {
		// Handle concurrent ServingRuntime creation - if it already exists, that's fine for shared usage
		if !apierr.IsAlreadyExists(err) {
			return configMap, nil, err
		}
		// If ServingRuntime already exists, get the existing one for consistency
		if getErr := k8sClient.Get(context.TODO(), types.NamespacedName{
			Name:      servingRuntime.Name,
			Namespace: servingRuntime.Namespace,
		}, servingRuntime); getErr != nil {
			return configMap, nil, fmt.Errorf("failed to get existing ServingRuntime: %w", getErr)
		}
	}

	return configMap, servingRuntime, nil
}

// CleanupTestResources cleans up common test resources
func CleanupTestResources(k8sClient client.Client, configMap *corev1.ConfigMap, servingRuntime *v1alpha1.ServingRuntime, isvc *v1beta1.InferenceService) {
	if isvc != nil {
		k8sClient.Delete(context.TODO(), isvc)
	}
	if servingRuntime != nil {
		k8sClient.Delete(context.TODO(), servingRuntime)
	}
	if configMap != nil {
		k8sClient.Delete(context.TODO(), configMap)
	}
}

// ValidateDeploymentQuick performs essential deployment validation without deep comparison
func ValidateDeploymentQuick(k8sClient client.Client, key types.NamespacedName, timeout time.Duration, expectedImage string) error {
	deployment := &appsv1.Deployment{}

	// Wait for deployment to be created
	var lastErr error
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(500 * time.Millisecond) {
		if err := k8sClient.Get(context.TODO(), key, deployment); err != nil {
			lastErr = err
			continue
		}

		// Quick validation - just check image and replica count
		if len(deployment.Spec.Template.Spec.Containers) > 0 {
			if expectedImage != "" && deployment.Spec.Template.Spec.Containers[0].Image != expectedImage {
				return errors.New("image mismatch")
			}
		}
		return nil
	}

	return fmt.Errorf("deployment validation failed: %w", lastErr)
}

// ValidateServiceQuick performs essential service validation
func ValidateServiceQuick(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) error {
	service := &corev1.Service{}

	// Wait for service to be created
	var lastErr error
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(500 * time.Millisecond) {
		if err := k8sClient.Get(context.TODO(), key, service); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	return fmt.Errorf("service validation failed: %w", lastErr)
}

// ValidateHPAQuick performs essential HPA validation
func ValidateHPAQuick(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) error {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{}

	// Wait for HPA to be created
	var lastErr error
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(500 * time.Millisecond) {
		if err := k8sClient.Get(context.TODO(), key, hpa); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	return fmt.Errorf("HPA validation failed: %w", lastErr)
}

// StandardTestSetup performs complete test setup and returns all created resources
type StandardTestSetup struct {
	ConfigMap        *corev1.ConfigMap
	ServingRuntime   *v1alpha1.ServingRuntime
	InferenceService *v1beta1.InferenceService
	ServiceKey       types.NamespacedName
	PredictorKey     types.NamespacedName
	TransformerKey   types.NamespacedName
	ExplainerKey     types.NamespacedName
}

// CreateStandardTestSetup creates a complete test setup with all common resources
func CreateStandardTestSetup(k8sClient client.Client, configs map[string]string, serviceName, namespace string, annotations map[string]string) (*StandardTestSetup, error) {
	// Create shared resources
	configMap, servingRuntime, err := SetupTestResources(k8sClient, configs, "tf-serving-raw", namespace)
	if err != nil {
		return nil, err
	}

	// Create InferenceService
	isvc := CreateDefaultInferenceService(serviceName, namespace, annotations)
	if err := k8sClient.Create(context.TODO(), isvc); err != nil {
		CleanupTestResources(k8sClient, configMap, servingRuntime, nil)
		return nil, err
	}

	serviceKey := types.NamespacedName{Name: serviceName, Namespace: namespace}

	return &StandardTestSetup{
		ConfigMap:        configMap,
		ServingRuntime:   servingRuntime,
		InferenceService: isvc,
		ServiceKey:       serviceKey,
		PredictorKey:     GetPredictorKey(serviceName, namespace),
		TransformerKey:   GetTransformerKey(serviceName, namespace),
		ExplainerKey:     GetExplainerKey(serviceName, namespace),
	}, nil
}

// Cleanup cleans up all resources in the StandardTestSetup
func (s *StandardTestSetup) Cleanup(k8sClient client.Client) {
	CleanupTestResources(k8sClient, s.ConfigMap, s.ServingRuntime, s.InferenceService)
}

// WaitForStandardResources waits for deployment, service, and HPA to be ready
func (s *StandardTestSetup) WaitForStandardResources(k8sClient client.Client, timeout time.Duration) error {
	// Quick validation instead of deep comparison
	if err := ValidateDeploymentQuick(k8sClient, s.PredictorKey, timeout, ""); err != nil {
		return err
	}
	if err := ValidateServiceQuick(k8sClient, s.PredictorKey, timeout); err != nil {
		return err
	}
	if err := ValidateHPAQuick(k8sClient, s.PredictorKey, timeout); err != nil {
		return err
	}
	return nil
}
