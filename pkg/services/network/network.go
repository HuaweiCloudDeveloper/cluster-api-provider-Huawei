package network

import (
	infrav1alpha1 "github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/api/v1alpha1"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

func (s *Service) ReconcileNetwork() error {
	klog.Infof("Reconciling network")

	// VPC
	if err := s.reconcileVPC(); err != nil {
		klog.Errorf("Failed to reconcile VPC: %v", err)
		conditions.MarkFalse(s.scope.InfraCluster(),
			infrav1alpha1.VpcReadyCondition,
			infrav1alpha1.VpcReconciliationFailedReason,
			clusterv1.ConditionSeverityError, "failed to reconcile VPC")
		return err
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1alpha1.VpcReadyCondition)

	// Subnets
	if err := s.reconcileSubnets(); err != nil {
		klog.Errorf("Failed to reconcile Subnets: %v", err)
		conditions.MarkFalse(
			s.scope.InfraCluster(),
			infrav1alpha1.SubnetsReadyCondition,
			infrav1alpha1.SubnetsReconciliationFailedReason,
			clusterv1.ConditionSeverityError, "failed to reconcile Subnets")
		return err
	}

	// TODO: Routing tables

	klog.Infof("Reconcile network completed successfully")
	return nil
}

func (s *Service) DeleteNetwork() error {
	klog.Infof("Deleting network")

	// Delete Subnets
	conditions.MarkFalse(
		s.scope.InfraCluster(),
		infrav1alpha1.SubnetsReadyCondition,
		clusterv1.DeletingReason,
		clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}
	if err := s.deleteSubnets(); err != nil {
		klog.Errorf("Failed to delete subnets: %v", err)
		conditions.MarkFalse(
			s.scope.InfraCluster(),
			infrav1alpha1.SubnetsReadyCondition,
			"DeletingFailed",
			clusterv1.ConditionSeverityWarning, "failed to delete subnets")
		return err
	}
	conditions.MarkFalse(
		s.scope.InfraCluster(),
		infrav1alpha1.SubnetsReadyCondition,
		clusterv1.DeletedReason,
		clusterv1.ConditionSeverityInfo, "")

	// TODO: Delete Route Tables

	// Delete VPC
	conditions.MarkFalse(
		s.scope.InfraCluster(),
		infrav1alpha1.VpcReadyCondition,
		clusterv1.DeletingReason,
		clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}
	if err := s.deleteVPC(); err != nil {
		klog.Errorf("Failed to delete VPC: %v", err)
		conditions.MarkFalse(
			s.scope.InfraCluster(),
			infrav1alpha1.VpcReadyCondition,
			"DeletingFailed",
			clusterv1.ConditionSeverityWarning, "failed to delete VPC")
		return err
	}
	conditions.MarkFalse(
		s.scope.InfraCluster(),
		infrav1alpha1.VpcReadyCondition,
		clusterv1.DeletedReason,
		clusterv1.ConditionSeverityInfo, "")

	klog.Infof("Delete network completed successfully")
	return nil
}
