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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	ClusterFinalizer = "huaweicloudcluster.infrastructure.cluster.x-k8s.io"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HuaweiCloudClusterSpec defines the desired state of HuaweiCloudCluster.
type HuaweiCloudClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NetworkSpec encapsulates the configuration options for HuaweiCloud network.
	NetworkSpec NetworkSpec `json:"network,omitempty"`

	// The ECS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// TODO, Network related fields need to be defined in the future
	// other fields may like S3, SSHKey, etc.
}

// HuaweiCloudClusterStatus defines the observed state of HuaweiCloudCluster.
type HuaweiCloudClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default=false
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HuaweiCloudCluster is the Schema for the huaweicloudclusters API.
type HuaweiCloudCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HuaweiCloudClusterSpec   `json:"spec,omitempty"`
	Status HuaweiCloudClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HuaweiCloudClusterList contains a list of HuaweiCloudCluster.
type HuaweiCloudClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiCloudCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HuaweiCloudCluster{}, &HuaweiCloudClusterList{})
}
