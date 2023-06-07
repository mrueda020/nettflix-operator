/*
Copyright 2023.

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

	webv1alpha1 "github.com/mrueda020/netflix-operator/api/v1alpha1"
	"github.com/mrueda020/netflix-operator/assets"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NetflixReconciler reconciles a Netflix object
type NetflixReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=web.example.com,resources=netflixes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=web.example.com,resources=netflixes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=web.example.com,resources=netflixes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Netflix object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *NetflixReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	netflixOperatorCR := &webv1alpha1.Netflix{}
	err := r.Get(ctx, req.NamespacedName, netflixOperatorCR)

	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("CR is not found for netflix operator")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Error getting operator resource object")
		return ctrl.Result{}, err
	}

	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, req.NamespacedName, deployment)

	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("No deployments for netflix, creating deployment")
			deployment.Namespace = req.Namespace
			deployment.Name = req.Name

			deploymentManifest, _ := assets.GetDeploymentFromFile("manifests/netflix_deployment.yaml")
			deploymentManifest.Spec.Replicas = pointer.Int32(3)

			err = r.Create(ctx, deploymentManifest)
			if err != nil {
				logger.Error(err, "Error creating netflix deployment.")
				return ctrl.Result{}, err
			}
		}
		logger.Error(err, "Error getting netflix deployment.")
		return ctrl.Result{}, err
	}

	service := corev1.Service{}
	err = r.Get(ctx, req.NamespacedName, &service)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("No deployments for netflix, creating deployment")
			service.Namespace = req.Namespace
			service.Name = req.Name

			serviceManifest, _ := assets.GetDeploymentFromFile("manifests/netflix_deployment.yaml")
			//serviceManifest.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *netflixOperatorCR.Spec.Port

			err = r.Create(ctx, serviceManifest)
			if err != nil {
				logger.Error(err, "Error creating netflix service.")
				return ctrl.Result{}, err
			}
		}
		logger.Error(err, "Error getting netflix service.")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetflixReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1alpha1.Netflix{}).
		Complete(r)
}
