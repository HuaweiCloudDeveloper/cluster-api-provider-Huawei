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
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util"
	capiannotations "sigs.k8s.io/cluster-api/util/annotations"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1alpha1 "github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/api/v1alpha1"
	"github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pkg/scope"
	"github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pkg/services/network"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/pkg/errors"
)

// HuaweiCloudClusterReconciler reconciles a HuaweiCloudCluster object
type HuaweiCloudClusterReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Credentials *basic.Credentials
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweicloudclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweicloudclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweicloudclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HuaweiCloudCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *HuaweiCloudClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := log.FromContext(ctx)

	// Fetch the HuaweiCloudCluster instance
	hcCluster := &infrav1alpha1.HuaweiCloudCluster{}
	err := r.Get(ctx, req.NamespacedName, hcCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// TODO: Set default values for HuaweiCloudCluster
	// hcCluster.Default()

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, hcCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}

	if cluster == nil {
		log.Info("Waiting for Cluster Controller to set OwnerRef on HuaweiCloudCluster")
		return reconcile.Result{}, nil
	}

	if capiannotations.IsPaused(cluster, hcCluster) {
		log.Info("HCCluster or linked Cluster is marked as paused. Won't reconcile")
		return reconcile.Result{}, nil
	}

	// Create the scope.
	clusterScope, err := scope.NewClusterScope(scope.ClusterScopeParams{
		Client:      r.Client,
		Logger:      &log,
		Cluster:     cluster,
		HCCluster:   hcCluster,
		Credentials: r.Credentials,
	})
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	defer func() {
		if err := clusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !hcCluster.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(clusterScope)
	}

	return r.reconcileNormal(clusterScope)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HuaweiCloudClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1alpha1.HuaweiCloudCluster{}).
		Named("huaweicloudcluster").
		Complete(r)
}

func (r *HuaweiCloudClusterReconciler) reconcileNormal(clusterScope *scope.ClusterScope) (reconcile.Result, error) {
	clusterScope.Logger.Info("Reconciling HuaweiCloudCluster")

	hccluster := clusterScope.HCCluster

	if controllerutil.AddFinalizer(hccluster, infrav1alpha1.ClusterFinalizer) {
		if err := clusterScope.PatchObject(); err != nil {
			return reconcile.Result{}, errors.Wrap(err, "failed to patch HuaweiCloudCluster with finalizer")
		}
	}

	// Reconcile network
	networkSvc, err := network.NewService(clusterScope)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to create network service")
	}
	if err := networkSvc.ReconcileNetwork(); err != nil {
		return reconcile.Result{RequeueAfter: 30 * time.Second}, errors.Wrap(err, "failed to reconcile network")
	}

	// TODO: Reconcile security group

	hccluster.Status.Ready = true
	return reconcile.Result{}, nil
}

func (r *HuaweiCloudClusterReconciler) reconcileDelete(clusterScope *scope.ClusterScope) error {
	// Reconcile network
	if !controllerutil.ContainsFinalizer(clusterScope.HCCluster, infrav1alpha1.ClusterFinalizer) {
		clusterScope.Logger.Info("No finalizer on HuaweiCloudCluster, skipping deletion reconciliation")
		return nil
	}

	clusterScope.Logger.Info("Deleting HuaweiCloudCluster")

	// Delete network
	networkSvc, err := network.NewService(clusterScope)
	if err != nil {
		return errors.Wrap(err, "failed to create network service")
	}
	if err := networkSvc.DeleteNetwork(); err != nil {
		return errors.Wrap(err, "failed to delete network")
	}

	// TODO: Delete security group

	controllerutil.RemoveFinalizer(clusterScope.HCCluster, infrav1alpha1.ClusterFinalizer)
	return nil
}
