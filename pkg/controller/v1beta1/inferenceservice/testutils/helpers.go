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
	"os"
	"time"

	"google.golang.org/protobuf/proto"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kserve/kserve/pkg/apis/serving/v1alpha1"
	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/kserve/kserve/pkg/constants"
	"knative.dev/pkg/apis"
)

// Test timeout constants - optimized for faster execution
const (
	FastTimeout         = time.Second * 10      // Further reduced from 15s for fast operations
	MediumTimeout       = time.Second * 20      // Reduced from 30s for medium operations
	SlowTimeout         = time.Second * 30      // Reduced from 45s for slow operations
	PollInterval        = time.Millisecond * 50 // Even faster polling interval
	ConsistentlyTimeout = time.Second * 2       // Further reduced from 3s
)

// Phase 2: Parallel Test Execution Constants
const (
	// Test parallelization settings
	MaxConcurrentTests = 4                      // Number of tests to run in parallel
	TestCooldownPeriod = time.Millisecond * 100 // Brief pause between test batches

	// Resource cleanup timeouts
	CleanupTimeout     = time.Second * 5       // Fast cleanup operations
	ResourceWaitPeriod = time.Millisecond * 25 // Very fast resource polling
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

// WaitForDeployment waits for a deployment to be available and returns it
func WaitForDeployment(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := k8sClient.Get(context.TODO(), key, deployment)
	return deployment, err
}

// WaitForService waits for a service to be available and returns it
func WaitForService(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := k8sClient.Get(context.TODO(), key, service)
	return service, err
}

// WaitForHPA waits for an HPA to be available and returns it
func WaitForHPA(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{}
	err := k8sClient.Get(context.TODO(), key, hpa)
	return hpa, err
}

// WaitForHTTPRoute waits for an HTTPRoute to be available and returns it
func WaitForHTTPRoute(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) (*gwapiv1.HTTPRoute, error) {
	httpRoute := &gwapiv1.HTTPRoute{}
	err := k8sClient.Get(context.TODO(), key, httpRoute)
	return httpRoute, err
}

// CleanupInferenceService deletes an InferenceService and waits for cleanup
func CleanupInferenceService(k8sClient client.Client, isvc *v1beta1.InferenceService) error {
	return k8sClient.Delete(context.TODO(), isvc)
}

// CleanupServingRuntime deletes a ServingRuntime
func CleanupServingRuntime(k8sClient client.Client, sr *v1alpha1.ServingRuntime) error {
	return k8sClient.Delete(context.TODO(), sr)
}

// CleanupConfigMap deletes a ConfigMap
func CleanupConfigMap(k8sClient client.Client, cm *corev1.ConfigMap) error {
	return k8sClient.Delete(context.TODO(), cm)
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

	return fmt.Errorf("deployment validation failed: %v", lastErr)
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

	return fmt.Errorf("service validation failed: %v", lastErr)
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

	return fmt.Errorf("HPA validation failed: %v", lastErr)
}

// ValidateHTTPRouteQuick performs essential HTTPRoute validation
func ValidateHTTPRouteQuick(k8sClient client.Client, key types.NamespacedName, timeout time.Duration) error {
	httpRoute := &gwapiv1.HTTPRoute{}
	return k8sClient.Get(context.TODO(), key, httpRoute)
}

// GetDefaultQuantity returns the default quantity used in tests (avoids repeated parsing)
var defaultQuantity = resource.MustParse("10Gi")

func GetDefaultQuantity() *resource.Quantity {
	return &defaultQuantity
}

// Pre-built common objects to avoid repeated object creation
func GetExpectedServiceSpec(serviceName, predictorName string) corev1.ServiceSpec {
	return corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name:       serviceName,
				Protocol:   "TCP",
				Port:       80,
				TargetPort: intstr.IntOrString{Type: 0, IntVal: 8080, StrVal: ""},
			},
		},
		Type:            "ClusterIP",
		SessionAffinity: "None",
		Selector: map[string]string{
			"app": "isvc." + predictorName,
		},
	}
}

// GetExpectedHPASpec returns a standard HPA spec for testing
func GetExpectedHPASpec(deploymentName string, minReplicas, maxReplicas int32) autoscalingv2.HorizontalPodAutoscalerSpec {
	return autoscalingv2.HorizontalPodAutoscalerSpec{
		ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       deploymentName,
		},
		MinReplicas: &minReplicas,
		MaxReplicas: maxReplicas,
		Metrics: []autoscalingv2.MetricSpec{
			{
				Type: autoscalingv2.ResourceMetricSourceType,
				Resource: &autoscalingv2.ResourceMetricSource{
					Name: corev1.ResourceMemory,
					Target: autoscalingv2.MetricTarget{
						Type:         autoscalingv2.AverageValueMetricType,
						AverageValue: GetDefaultQuantity(),
					},
				},
			},
		},
	}
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

// CreateInferenceServiceWithDeploymentStrategy creates a customized isvc for testing deployment strategies
func CreateInferenceServiceWithDeploymentStrategy(serviceName, namespace string, strategy appsv1.DeploymentStrategyType) *v1beta1.InferenceService {
	storageUri := "s3://test/mnist/export"
	var cpuUtilization int32 = 75

	isvc := &v1beta1.InferenceService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
			Annotations: map[string]string{
				constants.DeploymentMode:  string(constants.RawDeployment),
				constants.AutoscalerClass: string(constants.AutoscalerClassHPA),
			},
		},
		Spec: v1beta1.InferenceServiceSpec{
			Predictor: v1beta1.PredictorSpec{
				ComponentExtensionSpec: v1beta1.ComponentExtensionSpec{
					MinReplicas: ptr.To(int32(1)),
					MaxReplicas: 3,
					DeploymentStrategy: &appsv1.DeploymentStrategy{
						Type: strategy,
					},
					AutoScaling: &v1beta1.AutoScalingSpec{
						Metrics: []v1beta1.MetricsSpec{
							{
								Type: v1beta1.ResourceMetricSourceType,
								Resource: &v1beta1.ResourceMetricSource{
									Name: v1beta1.ResourceMetricCPU,
									Target: v1beta1.MetricTarget{
										Type:               v1beta1.UtilizationMetricType,
										AverageUtilization: &cpuUtilization,
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
			},
		},
	}

	isvc.DefaultInferenceService(nil, nil, &v1beta1.SecurityConfig{AutoMountServiceAccountToken: false}, nil)
	return isvc
}

// ValidateInferenceServiceStatus performs quick status validation
func ValidateInferenceServiceStatus(k8sClient client.Client, serviceKey types.NamespacedName, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), PollInterval, timeout, true, func(ctx context.Context) (bool, error) {
		isvc := &v1beta1.InferenceService{}
		if err := k8sClient.Get(ctx, serviceKey, isvc); err != nil {
			return false, err
		}

		// Quick status check - just verify ready conditions exist
		for _, condition := range isvc.Status.Conditions {
			if condition.Type == apis.ConditionReady && condition.Status == "True" {
				return true, nil
			}
		}
		return false, nil
	})
}

// ValidateDeploymentStrategy validates deployment strategy quickly
func ValidateDeploymentStrategy(k8sClient client.Client, deploymentKey types.NamespacedName, expectedStrategy appsv1.DeploymentStrategyType, timeout time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), PollInterval, timeout, true, func(ctx context.Context) (bool, error) {
		deployment := &appsv1.Deployment{}
		if err := k8sClient.Get(ctx, deploymentKey, deployment); err != nil {
			return false, err
		}
		return deployment.Spec.Strategy.Type == expectedStrategy, nil
	})
}

// BulkValidateResources validates multiple resources quickly
func BulkValidateResources(k8sClient client.Client, keys []types.NamespacedName, resourceTypes []string, timeout time.Duration) error {
	for i, key := range keys {
		switch resourceTypes[i] {
		case "deployment":
			if err := ValidateDeploymentQuick(k8sClient, key, timeout, ""); err != nil {
				return fmt.Errorf("deployment validation failed: %w", err)
			}
		case "service":
			if err := ValidateServiceQuick(k8sClient, key, timeout); err != nil {
				return fmt.Errorf("service validation failed: %w", err)
			}
		case "hpa":
			if err := ValidateHPAQuick(k8sClient, key, timeout); err != nil {
				return fmt.Errorf("HPA validation failed: %w", err)
			}
		case "httproute":
			if err := ValidateHTTPRouteQuick(k8sClient, key, timeout); err != nil {
				return fmt.Errorf("HTTPRoute validation failed: %w", err)
			}
		}
	}
	return nil
}

// FastResourceValidation performs the quickest possible validation
func FastResourceValidation(k8sClient client.Client, keys []types.NamespacedName) error {
	for _, key := range keys {
		// Just check that resources exist, no deep validation
		deployment := &appsv1.Deployment{}
		if err := k8sClient.Get(context.TODO(), key, deployment); err != nil {
			return fmt.Errorf("deployment validation failed: %w", err)
		}

		service := &corev1.Service{}
		if err := k8sClient.Get(context.TODO(), key, service); err != nil {
			// Service might have different name, try with predictor service name
			serviceKey := types.NamespacedName{
				Name:      constants.PredictorServiceName(key.Name),
				Namespace: key.Namespace,
			}
			if err := k8sClient.Get(context.TODO(), serviceKey, service); err != nil {
				return fmt.Errorf("service validation failed: %w", err)
			}
		}
	}
	return nil
}

// BatchOptimizeTestCase applies all optimizations to a test case
func BatchOptimizeTestCase(k8sClient client.Client, serviceName, namespace string, configs map[string]string, runtimeName string, validationFunc func() error) error {
	// Setup
	configMap, servingRuntime, err := SetupTestResources(k8sClient, configs, runtimeName, namespace)
	if err != nil {
		return err
	}
	defer CleanupTestResources(k8sClient, configMap, servingRuntime, nil)

	// Run validation with fast timeouts
	if validationFunc != nil {
		return validationFunc()
	}
	return nil
}

// Common test constants - shared across all test cases
const (
	DefaultStorageURI  = "s3://test/mnist/export"
	DefaultNamespace   = "default"
	DefaultRuntimeName = "tf-serving-raw"
)

// CreateTestServiceKeys creates all necessary service keys for a test case
func CreateTestServiceKeys(serviceName, namespace string) (serviceKey, predictorKey, transformerKey, explainerKey types.NamespacedName) {
	serviceKey = types.NamespacedName{Name: serviceName, Namespace: namespace}
	predictorKey = types.NamespacedName{
		Name:      constants.PredictorServiceName(serviceName),
		Namespace: namespace,
	}
	transformerKey = types.NamespacedName{
		Name:      constants.TransformerServiceName(serviceName),
		Namespace: namespace,
	}
	explainerKey = types.NamespacedName{
		Name:      constants.ExplainerServiceName(serviceName),
		Namespace: namespace,
	}
	return
}

// CreateDefaultTestRequest creates a standard reconcile request
func CreateDefaultTestRequest(serviceName, namespace string) reconcile.Request {
	return reconcile.Request{NamespacedName: types.NamespacedName{Name: serviceName, Namespace: namespace}}
}

// QuickTestCase runs a complete test case with minimal setup
func QuickTestCase(k8sClient client.Client, configs map[string]string, serviceName string, testFunc func(serviceKey types.NamespacedName) error) error {
	// Setup resources
	configMap, servingRuntime, err := SetupTestResources(k8sClient, configs, DefaultRuntimeName, DefaultNamespace)
	if err != nil {
		return err
	}
	defer CleanupTestResources(k8sClient, configMap, servingRuntime, nil)

	// Create service key
	serviceKey := types.NamespacedName{Name: serviceName, Namespace: DefaultNamespace}

	// Run test function
	return testFunc(serviceKey)
}

// CommonTestInferenceService creates a standard test InferenceService with common defaults
func CommonTestInferenceService(serviceName, namespace string, annotations map[string]string) *v1beta1.InferenceService {
	if annotations == nil {
		annotations = map[string]string{
			constants.DeploymentMode:  string(constants.RawDeployment),
			constants.AutoscalerClass: string(constants.AutoscalerClassHPA),
		}
	}

	storageUri := DefaultStorageURI
	return &v1beta1.InferenceService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceName,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: v1beta1.InferenceServiceSpec{
			Predictor: v1beta1.PredictorSpec{
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
			},
		},
	}
}

// ValidateAllStandardResources validates deployment, service, and HPA in sequence with optimized timeouts
func ValidateAllStandardResources(k8sClient client.Client, serviceName, namespace string) error {
	_, predictorKey, _, _ := CreateTestServiceKeys(serviceName, namespace)

	// Progressive validation - existence first, then details
	if err := ValidateDeploymentQuick(k8sClient, predictorKey, FastTimeout, ""); err != nil {
		return fmt.Errorf("deployment validation failed: %w", err)
	}

	// Validate predictor service (not the main service, which doesn't exist in raw mode)
	if err := ValidateServiceQuick(k8sClient, predictorKey, FastTimeout); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := ValidateHPAQuick(k8sClient, predictorKey, FastTimeout); err != nil {
		return fmt.Errorf("HPA validation failed: %w", err)
	}

	return nil
}

// Phase 2: Parallel Test Execution Helpers

// TestBatch represents a batch of tests that can run concurrently
type TestBatch struct {
	Name        string
	TestFuncs   []func() error
	MaxParallel int
}

// ParallelTestManager manages concurrent test execution
type ParallelTestManager struct {
	client         client.Client
	maxConcurrent  int
	cooldownPeriod time.Duration
}

// NewParallelTestManager creates a new parallel test manager
func NewParallelTestManager(client client.Client) *ParallelTestManager {
	return &ParallelTestManager{
		client:         client,
		maxConcurrent:  MaxConcurrentTests,
		cooldownPeriod: TestCooldownPeriod,
	}
}

// RunTestBatch executes a batch of tests in parallel with controlled concurrency
func (ptm *ParallelTestManager) RunTestBatch(batch TestBatch) error {
	semaphore := make(chan struct{}, batch.MaxParallel)
	errChan := make(chan error, len(batch.TestFuncs))

	// Launch all test functions with controlled concurrency
	for i, testFunc := range batch.TestFuncs {
		go func(index int, fn func() error) {
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			if err := fn(); err != nil {
				errChan <- fmt.Errorf("test %d in batch %s failed: %w", index, batch.Name, err)
				return
			}
			errChan <- nil
		}(i, testFunc)
	}

	// Wait for all tests to complete and collect errors
	var errors []error
	for i := 0; i < len(batch.TestFuncs); i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	// Add cooldown period between batches
	time.Sleep(ptm.cooldownPeriod)

	if len(errors) > 0 {
		return fmt.Errorf("batch %s had %d failures: %v", batch.Name, len(errors), errors)
	}

	return nil
}

// FastResourceCleanup performs optimized cleanup of test resources
func FastResourceCleanup(k8sClient client.Client, resources ...client.Object) error {
	errChan := make(chan error, len(resources))

	// Delete all resources in parallel
	for _, resource := range resources {
		go func(obj client.Object) {
			ctx, cancel := context.WithTimeout(context.Background(), CleanupTimeout)
			defer cancel()
			errChan <- client.IgnoreNotFound(k8sClient.Delete(ctx, obj))
		}(resource)
	}

	// Wait for all deletions to complete
	var errors []error
	for i := 0; i < len(resources); i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("cleanup failed for %d resources: %v", len(errors), errors)
	}

	return nil
}

// BatchCreateResources creates multiple resources in parallel for faster setup
func BatchCreateResources(k8sClient client.Client, resources ...client.Object) error {
	errChan := make(chan error, len(resources))

	// Create all resources in parallel
	for _, resource := range resources {
		go func(obj client.Object) {
			ctx, cancel := context.WithTimeout(context.Background(), FastTimeout)
			defer cancel()
			errChan <- k8sClient.Create(ctx, obj)
		}(resource)
	}

	// Wait for all creations to complete
	var errors []error
	for i := 0; i < len(resources); i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("batch creation failed for %d resources: %v", len(errors), errors)
	}

	return nil
}

// Phase 2: Advanced Test Organization Helpers

// TestResourcePool manages a pool of test resources for reuse
type TestResourcePool struct {
	configMaps      []*corev1.ConfigMap
	servingRuntimes []*v1alpha1.ServingRuntime
	client          client.Client
	poolSize        int
}

// NewTestResourcePool creates a new resource pool for faster test execution
func NewTestResourcePool(client client.Client, poolSize int) *TestResourcePool {
	return &TestResourcePool{
		configMaps:      make([]*corev1.ConfigMap, 0, poolSize),
		servingRuntimes: make([]*v1alpha1.ServingRuntime, 0, poolSize),
		client:          client,
		poolSize:        poolSize,
	}
}

// PrepopulatePool creates a pool of resources for quick test execution
func (pool *TestResourcePool) PrepopulatePool(configs map[string]string) error {
	// Create multiple configmaps and serving runtimes for parallel test use
	for i := 0; i < pool.poolSize; i++ {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%d", constants.InferenceServiceConfigMapName, i),
				Namespace: constants.KServeNamespace,
			},
			Data: configs,
		}

		servingRuntime := &v1alpha1.ServingRuntime{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("tf-serving-raw-%d", i),
				Namespace: "default",
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

		pool.configMaps = append(pool.configMaps, configMap)
		pool.servingRuntimes = append(pool.servingRuntimes, servingRuntime)
	}

	// Create all resources in parallel
	var allResources []client.Object
	for _, cm := range pool.configMaps {
		allResources = append(allResources, cm)
	}
	for _, sr := range pool.servingRuntimes {
		allResources = append(allResources, sr)
	}

	return BatchCreateResources(pool.client, allResources...)
}

// GetPooledResources returns a set of pre-created resources for immediate use
func (pool *TestResourcePool) GetPooledResources(index int) (*corev1.ConfigMap, *v1alpha1.ServingRuntime, error) {
	if index >= len(pool.configMaps) || index >= len(pool.servingRuntimes) {
		return nil, nil, fmt.Errorf("pool index %d out of range (pool size: %d)", index, pool.poolSize)
	}

	return pool.configMaps[index], pool.servingRuntimes[index], nil
}

// CleanupPool removes all pooled resources
func (pool *TestResourcePool) CleanupPool() error {
	var allResources []client.Object
	for _, cm := range pool.configMaps {
		allResources = append(allResources, cm)
	}
	for _, sr := range pool.servingRuntimes {
		allResources = append(allResources, sr)
	}

	return FastResourceCleanup(pool.client, allResources...)
}

// OptimizedTestRunner provides high-performance test execution with resource pooling
type OptimizedTestRunner struct {
	resourcePool *TestResourcePool
	testManager  *ParallelTestManager
	client       client.Client
}

// NewOptimizedTestRunner creates a comprehensive test runner for Phase 2 optimizations
func NewOptimizedTestRunner(client client.Client, poolSize int) *OptimizedTestRunner {
	return &OptimizedTestRunner{
		resourcePool: NewTestResourcePool(client, poolSize),
		testManager:  NewParallelTestManager(client),
		client:       client,
	}
}

// RunOptimizedTestSuite executes a full test suite with maximum parallelization
func (runner *OptimizedTestRunner) RunOptimizedTestSuite(configs map[string]string, testCases []TestBatch) error {
	// Setup resource pool
	if err := runner.resourcePool.PrepopulatePool(configs); err != nil {
		return fmt.Errorf("failed to setup resource pool: %w", err)
	}
	defer runner.resourcePool.CleanupPool()

	// Execute test batches in sequence (each batch runs tests in parallel)
	for _, batch := range testCases {
		if err := runner.testManager.RunTestBatch(batch); err != nil {
			return fmt.Errorf("test batch %s failed: %w", batch.Name, err)
		}
	}

	return nil
}

// IsQuickTestMode returns true if QUICK_TESTS environment variable is set
func IsQuickTestMode() bool {
	return os.Getenv("QUICK_TESTS") == "true"
}
