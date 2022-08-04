/*
Copyright 2022.

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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "github.com/jeantsai/sample-operator/api/v1alpha1"
	"github.com/jeantsai/sample-operator/assets"
)

// LegacyAppReconciler reconciles a LegacyApp object
type LegacyAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operator.jeantsai.cn,resources=legacyapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.jeantsai.cn,resources=legacyapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.jeantsai.cn,resources=legacyapps/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LegacyApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *LegacyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	cr := &operatorv1alpha1.LegacyApp{}
	err := r.Get(ctx, req.NamespacedName, cr)
	if err != nil && errors.IsNotFound(err) {
		logger.Error(err, "Failed to find operator/legacyapp")
		return ctrl.Result{}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get custom resource of operator/legacyapp")
		return ctrl.Result{}, err
	}

	dep := &appsv1.Deployment{}
	create := false
	err = r.Get(ctx, req.NamespacedName, dep)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating deployment ...")
		create = true
		dep = assets.GetLegacyDeploymentFromEmbeddedFile("manifests/legacy-deployment.yaml")
		logger.Info("Initializing deployment: %+v", "deployment:", dep)
	} else if err != nil {
		logger.Error(err, "Failed to get existing deployment of legacyapp")
		return ctrl.Result{}, err
	}

	dep.Namespace = req.Namespace
	dep.Name = req.Name

	if cr.Spec.Replicas != nil {
		dep.Spec.Replicas = cr.Spec.Replicas
	}

	if cr.Spec.Port != nil {
		dep.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *cr.Spec.Port
	}

	if create {
		err = r.Create(ctx, dep)
	} else {
		err = r.Update(ctx, dep)
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *LegacyAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.LegacyApp{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
