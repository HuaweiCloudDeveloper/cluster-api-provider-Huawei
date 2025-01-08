package network

import "k8s.io/klog/v2"

func (s *Service) ReconcileNetwork() error {
	klog.Infof("Reconciling network")

	// VPC
	if err := s.reconcileVPC(); err != nil {
		klog.Errorf("Failed to reconcile VPC: %v", err)
		return err
	}

	// TODO: Subnets

	// TODO: Public IPs

	// TODO: Routing tables

	klog.Infof("Reconcile network completed successfully")
	return nil
}

func (s *Service) DeleteNetwork() error {
	klog.Infof("Deleting network")

	// TODO: Delete Subnets

	// TODO: Delete Public IPs

	// TODO: Delete Route Tables

	// Delete VPC
	if err := s.deleteVPC(); err != nil {
		klog.Errorf("Failed to delete VPC: %v", err)
		return err
	}

	klog.Infof("Delete network completed successfully")
	return nil
}
