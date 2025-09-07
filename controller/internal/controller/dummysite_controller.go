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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	stabledwkv1 "stable.dwk/api/v1"
)

// DummySiteReconciler reconciles a DummySite object
type DummySiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stable.dwk.stable.dwk,resources=dummysites/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DummySite object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *DummySiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO(user): your logic here
	var dummysite stabledwkv1.DummySite
	if err := r.Get(ctx, req.NamespacedName, &dummysite); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      dummysite.Name + "-pod",
			Namespace: dummysite.Namespace,
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "dummysite-sample",
					Image: "aritrabb/dummysite:latest",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "TARGET_URL",
							Value: dummysite.Spec.WebsiteUrl,
						},
					},
				},
			},
		},
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, pod, func() error {
		return ctrl.SetControllerReference(&dummysite, pod, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, err
	}
	fmt.Printf("Pod deployed successfully\n")

	dummysite.Status.Conditions = []v1.Condition{
		{
			Type:               "Available",
			Status:             v1.ConditionTrue,
			Reason:             "Deployed",
			Message:            "Nginx serving fetched website",
			LastTransitionTime: v1.Now(),
		},
	}

	if err := r.Status().Update(ctx, &dummysite); err != nil {
		fmt.Println("failed to update status")
		return ctrl.Result{}, err
	}

	fmt.Printf("All good and green in reconciliation!")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummySiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stabledwkv1.DummySite{}).
		Named("dummysite").
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
