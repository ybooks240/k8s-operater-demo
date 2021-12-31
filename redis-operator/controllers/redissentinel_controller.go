/*
Copyright 2021 James.Liu.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	redisv1 "github.com/ybooks240/redisSentinel/api/v1"
)

// RedisSentinelReconciler reconciles a RedisSentinel object
type RedisSentinelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=redis.buleye.com,resources=redissentinels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=redis.buleye.com,resources=redissentinels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=redis.buleye.com,resources=redissentinels/finalizers,verbs=update
//+kubebuilder:rbac:groups=corev1,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=appsv1,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RedisSentinel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *RedisSentinelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog := log.FromContext(ctx)
	var redisSentinel redisv1.RedisSentinel

	klog.Info("正在调谐...\n")
	// 获取对象
	if err := r.Get(ctx, req.NamespacedName, &redisSentinel); err != nil {
		klog.WithName("IgnoreNotFound")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// svc 调谐
	for _, v := range serviceType {
		var redisSvc corev1.Service
		redisSvc.Name = redisSentinel.GetRedisSentinelName(v)
		redisSvc.Namespace = redisSentinel.Namespace

		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			re, err := ctrl.CreateOrUpdate(ctx, r.Client, &redisSvc, func() error {
				MutateHeadlessSvc(&redisSentinel, &redisSvc, v)
				return controllerutil.SetControllerReference(&redisSentinel, &redisSvc, r.Scheme)
			})
			klog.Info("CreateORUPtdae", "service", re)
			return err
		}); err != nil {
			return ctrl.Result{}, err
		}
	}
	// pod 调谐
	for _, v := range serviceType {
		var sts appsv1.StatefulSet
		sts.Name = redisSentinel.GetRedisSentinelName(v)
		fmt.Println(sts.Name)
		sts.Namespace = redisSentinel.Namespace

		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			re, err := ctrl.CreateOrUpdate(ctx, r.Client, &sts, func() error {
				MutateStatefulSet(&redisSentinel, &sts, v)
				return controllerutil.SetControllerReference(&redisSentinel, &sts, r.Scheme)
			})
			klog.Info("CreateORUPtdae", "StatefulSet", re)
			return err
		}); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RedisSentinelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&redisv1.RedisSentinel{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
