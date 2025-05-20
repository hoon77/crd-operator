/*
Copyright 2025.

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

package controller

import (
	"context"
	webappv1 "github.com/hoon77/crd-operator/api/v1"
	"github.com/hoon77/crd-operator/internal/pkg/resources"
	"github.com/hoon77/crd-operator/internal/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// WebAppReconciler reconciles a WebApp object
type WebAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.crdlego.com,resources=webapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.crdlego.com,resources=webapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.crdlego.com,resources=webapps/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *WebAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	klog.Info("Call Reconcile", "NamespacedName", req.NamespacedName)
	// TODO(user): your logic here

	var webapp webappv1.WebApp
	if err := r.Get(ctx, req.NamespacedName, &webapp); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// detect webapp deletion
	if !webapp.DeletionTimestamp.IsZero() {
		klog.Infof("Webapp %s/%s is being deleted. Cleaning up...", webapp.Namespace, webapp.Name)

		// delete deployment
		_ = r.Delete(ctx, &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
			Namespace: webapp.Namespace,
			Name:      webapp.Name,
		}})
		// delete service
		_ = r.Delete(ctx, &corev1.Service{ObjectMeta: metav1.ObjectMeta{
			Namespace: webapp.Namespace,
			Name:      webapp.Name,
		}})
		// delete configmap
		_ = r.Delete(ctx, &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Namespace: webapp.Namespace,
			Name:      webapp.Name + resources.ConfigMapSuffix,
		}})

		// delete ingress
		_ = r.Delete(ctx, &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
			Namespace: webapp.Namespace,
			Name:      webapp.Name + resources.IngressTLSSecretNameMaSuffix,
		}})

		// remove finalizer
		controllerutil.RemoveFinalizer(&webapp, resources.WebAppFinalizer)
		if err := r.Update(ctx, &webapp); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// add finalizer to webapp
	if !controllerutil.ContainsFinalizer(&webapp, resources.WebAppFinalizer) {
		controllerutil.AddFinalizer(&webapp, resources.WebAppFinalizer)
		if err := r.Update(ctx, &webapp); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Create configmap
	createConfigmap := resources.BuildConfigMap(&webapp)
	if err := utils.SetOwnerRefence(&webapp, createConfigmap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	foundConfigmap := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Namespace: createConfigmap.Namespace, Name: createConfigmap.Name}, foundConfigmap)
	if err != nil && errors.IsNotFound(err) {
		if err := r.Create(ctx, createConfigmap); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		if !reflect.DeepEqual(foundConfigmap.Data, createConfigmap.Data) {
			log.Info("New Configmap Data", "Data", createConfigmap.Data)
			foundConfigmap.Data = createConfigmap.Data
			if err := r.Update(ctx, foundConfigmap); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Create deployment
	createDeploy := resources.BuildDeployment(&webapp)
	if err := utils.SetOwnerRefence(&webapp, createDeploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	foundDeploy := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Namespace: webapp.Namespace, Name: webapp.Name}, foundDeploy)
	if err != nil && errors.IsNotFound(err) {
		if err := r.Create(ctx, createDeploy); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// compare config-hash
		oldHash := foundDeploy.Spec.Template.Annotations[resources.WebAppHashKey]
		newHash := createDeploy.Spec.Template.Annotations[resources.WebAppHashKey]
		if oldHash != newHash {
			foundDeploy.Spec.Template.Annotations = createDeploy.Spec.Template.Annotations
			if err := r.Update(ctx, foundDeploy); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Create service
	createSvc := resources.BuildService(&webapp)
	if err := utils.SetOwnerRefence(&webapp, createSvc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	foundSvc := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: webapp.Namespace, Name: webapp.Name}, foundSvc); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, createSvc); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if webapp.Spec.Ingress != nil && webapp.Spec.Ingress.Enabled {
		createIngress := resources.BuildIngress(&webapp)
		if err := utils.SetOwnerRefence(&webapp, createIngress, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		var foundIngress networkingv1.Ingress
		if err := r.Get(ctx, types.NamespacedName{Namespace: webapp.Namespace, Name: webapp.Name}, &foundIngress); err != nil {
			if errors.IsNotFound(err) {
				err = r.Create(ctx, createIngress)
			}
		} else if err == nil {
			createIngress.ResourceVersion = foundIngress.ResourceVersion
			if err = r.Update(ctx, createIngress); err != nil {
				return ctrl.Result{}, err
			}
		}

	}

	// Get deployment status availableReplicas
	availableReplicas := foundDeploy.Status.AvailableReplicas
	if webapp.Status.AvailableReplicas != availableReplicas {
		webapp.Status.AvailableReplicas = availableReplicas
		if err := r.Status().Update(ctx, &webapp); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.WebApp{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&networkingv1.Ingress{}).
		Named("webapp").
		Complete(r)
}
