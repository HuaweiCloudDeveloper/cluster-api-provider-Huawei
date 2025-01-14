/*
Copyright 2022 The Kubernetes Authors.

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

package v1alpha1

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	// VpcReadyCondition reports on the successful reconciliation of a VPC.
	VpcReadyCondition clusterv1.ConditionType = "VpcReady"
	// VpcCreationStartedReason used when attempting to create a VPC for a managed cluster.
	// Will not be applied to unmanaged clusters.
	VpcCreationStartedReason = "VpcCreationStarted"
	// VpcReconciliationFailedReason used when errors occur during VPC reconciliation.
	VpcReconciliationFailedReason = "VpcReconciliationFailed"
)

const (
	// SubnetsReadyCondition reports on the successful reconciliation of subnets.
	SubnetsReadyCondition clusterv1.ConditionType = "SubnetsReady"
	// SubnetsReconciliationFailedReason used to report failures while reconciling subnets.
	SubnetsReconciliationFailedReason = "SubnetsReconciliationFailed"
)

const (
	// ClusterSecurityGroupsReadyCondition reports successful reconciliation of security groups.
	ClusterSecurityGroupsReadyCondition clusterv1.ConditionType = "ClusterSecurityGroupsReady"
	// ClusterSecurityGroupReconciliationFailedReason used when any errors occur during reconciliation of security groups.
	ClusterSecurityGroupReconciliationFailedReason = "SecurityGroupReconciliationFailed"
)
