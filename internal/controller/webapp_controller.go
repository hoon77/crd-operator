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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

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
	log.Info("Call Reconcile", "NamespacedName", req.NamespacedName)
	// TODO(user): your logic here

	var webapp webappv1.WebApp
	if err := r.Get(ctx, req.NamespacedName, &webapp); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
	}

	// Create deployment
	createDeploy := resources.BuildDeployment(&webapp)
	if err := utils.SetOwnerRefence(&webapp, createDeploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	foundDeploy := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: webapp.Namespace, Name: webapp.Name}, foundDeploy); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, createDeploy); err != nil {
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
		Named("webapp").
		Complete(r)
}
