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

package autoscaler

import (
	"context"
	"fmt"

	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/kserve/kserve/pkg/constants"
	hpa "github.com/kserve/kserve/pkg/controller/v1beta1/inferenceservice/reconcilers/hpa"
	"github.com/kserve/kserve/pkg/controller/v1beta1/inferenceservice/reconcilers/keda"
)

var log = logf.Log.WithName("ExternalAutoscaler")

// Autoscaler Interface implemented by all autoscalers
type Autoscaler interface {
	Reconcile(ctx context.Context) error
	SetControllerReferences(owner metav1.Object, scheme *runtime.Scheme) error
}

// AutoscalerReconciler is the struct of Raw K8S Object
type AutoscalerReconciler struct {
	client       client.Client
	scheme       *runtime.Scheme
	Autoscaler   Autoscaler
	componentExt *v1beta1.ComponentExtensionSpec
}

func NewAutoscalerReconciler(client client.Client,
	scheme *runtime.Scheme,
	componentMeta metav1.ObjectMeta,
	componentExt *v1beta1.ComponentExtensionSpec,
	configMap *corev1.ConfigMap,
) (*AutoscalerReconciler, error) {
	as, err := createAutoscaler(client, scheme, componentMeta, componentExt, configMap)
	if err != nil {
		return nil, err
	}
	return &AutoscalerReconciler{
		client:       client,
		scheme:       scheme,
		Autoscaler:   as,
		componentExt: componentExt,
	}, err
}

func getAutoscalerClass(metadata metav1.ObjectMeta) constants.AutoscalerClassType {
	annotations := metadata.Annotations
	if value, ok := annotations[constants.AutoscalerClass]; ok {
		return constants.AutoscalerClassType(value)
	} else {
		return constants.DefaultAutoscalerClass
	}
}

func createAutoscaler(client client.Client,
	scheme *runtime.Scheme, componentMeta metav1.ObjectMeta,
	componentExt *v1beta1.ComponentExtensionSpec,
	configMap *corev1.ConfigMap,
) (Autoscaler, error) {
	ac := getAutoscalerClass(componentMeta)
	switch ac {
	case constants.AutoscalerClassHPA:
		return hpa.NewHPAReconciler(client, scheme, componentMeta, componentExt)
	case constants.AutoscalerClassKeda:
		return keda.NewKedaReconciler(client, scheme, componentMeta, componentExt, configMap)
	case constants.AutoscalerClassExternal:
		return &NoOpAutoscaler{client: client}, nil
	default:
		return nil, fmt.Errorf("unknown autoscaler class type: %v", ac)
	}
}

// Reconcile autoscaling resources for HPA, KEDA ScaledObject.
func (r *AutoscalerReconciler) Reconcile(ctx context.Context) error {
	// reconcile Autoscaling resources
	err := r.Autoscaler.Reconcile(ctx)
	if err != nil {
		return err
	}
	return nil
}

// NoOpAutoscaler Autoscaler that does nothing. Can be used to disable creation of autoscaler resources.
type NoOpAutoscaler struct {
	client client.Client
}

func (n *NoOpAutoscaler) Reconcile(ctx context.Context) error {
	// NoOp autoscaler doesn't create any resources, but cleans up any existing ones
	// that might have been created by other autoscalers previously
	n.CleanupHPA(ctx)
	n.CleanupKeda(ctx)
	return nil
}

// CleanupHPA removes any associated HPAs for isvc
func (n *NoOpAutoscaler) CleanupHPA(ctx context.Context) error {
	// Check if HPA exists
	hpaList := &autoscalingv2.HorizontalPodAutoscalerList{}
	if err := n.client.List(ctx, hpaList); err != nil {
		log.Error(err, "Failed to list HPAs")
		return err
	}

	for _, hpaObj := range hpaList.Items {
		// Check if this HPA has the annotation that associates it with KServe
		if value, ok := hpaObj.Annotations[constants.AutoscalerClass]; ok &&
			constants.AutoscalerClassType(value) == constants.AutoscalerClassHPA {
			log.Info("Deleting HPA", "name", hpaObj.Name, "namespace", hpaObj.Namespace)
			if err := n.client.Delete(ctx, &hpaObj); err != nil {
				log.Error(err, "Failed to delete HPA", "name", hpaObj.Name, "namespace", hpaObj.Namespace)
				return err
			}
		}
	}
	return nil
}

// CleanupKeda removes any associated KEDA ScaledObjects for isvc
func (n *NoOpAutoscaler) CleanupKeda(ctx context.Context) error {
	// Check if KEDA ScaledObject CRD exists before attempting cleanup
	_, err := n.client.RESTMapper().RESTMapping(kedav1alpha1.GroupVersion.WithKind(constants.KedaScaledObjectKind).GroupKind())
	if err != nil {
		return nil
	} else {
		// Check if KEDA ScaledObject exists
		kedaList := &kedav1alpha1.ScaledObjectList{}
		if err := n.client.List(ctx, kedaList); err != nil {
			log.Error(err, "Failed to list KEDA ScaledObjects")
			return err
		}

		for _, kedaObj := range kedaList.Items {
			// Check if this ScaledObject has the annotation that associates it with KServe
			if value, ok := kedaObj.Annotations[constants.AutoscalerClass]; ok &&
				constants.AutoscalerClassType(value) == constants.AutoscalerClassKeda {
				log.Info("Deleting KEDA ScaledObject", "name", kedaObj.Name, "namespace", kedaObj.Namespace)
				if err := n.client.Delete(ctx, &kedaObj); err != nil {
					log.Error(err, "Failed to delete KEDA ScaledObject", "name", kedaObj.Name, "namespace", kedaObj.Namespace)
					return err
				}
			}
		}
	}
	return nil
}

func (n *NoOpAutoscaler) SetControllerReferences(owner metav1.Object, scheme *runtime.Scheme) error {
	return nil
}
