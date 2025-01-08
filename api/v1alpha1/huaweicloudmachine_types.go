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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HuaweiCloudMachineSpec defines the desired state of HuaweiCloudMachine.
type HuaweiCloudMachineSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// 弹性云服务器uuid。
	// ProviderID is the unique identifier as specified by the cloud provider.
	ProviderID *string `json:"providerID,omitempty"`

	// InstanceID is the ECS instance ID for this machine.
	InstanceID *string `json:"instanceID,omitempty"`

	// 镜像ID或者镜像资源的URL
	// ImageRef is the reference from which to create the machine instance.
	ImageRef *string `json:"imageRef,omitempty"`

	// FlavorRef is similar to instanceType.
	// FlavorRef is the type of instance to create. Example: s2.small.1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=2
	FlavorRef string `json:"flavorRef"`

	// TODO, need to define the type of Volume struct in the future
	// RootVolume encapsulates the configuration options for the root volume
	// +optional
	// RootVolume *model.RootVolume `json:"rootVolume,omitempty"`
	// Configuration options for the data storage volumes.
	// +optional
	// DataVolumes *[]model.DataVolumes `json:"dataVolumes,omitempty"`

	// TODO, more fields need to be defined in the future
	// AdminPass *string `json:"admin_pass,omitempty"`
	// NetConfig *NetConfig `json:"net_config"`
	// Bandwidth *BandwidthConfig `json:"bandwidth,omitempty"`
}

// HuaweiCloudMachineStatus defines the observed state of HuaweiCloudMachine.
type HuaweiCloudMachineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the ECS instance associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// InstanceState is the state of the ECS instance for this machine.
	// +optional
	InstanceState *InstanceState `json:"instanceState,omitempty"`

	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// TODO add conditions and more fields in the future
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HuaweiCloudMachine is the Schema for the huaweicloudmachines API.
type HuaweiCloudMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HuaweiCloudMachineSpec   `json:"spec,omitempty"`
	Status HuaweiCloudMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HuaweiCloudMachineList contains a list of HuaweiCloudMachine.
type HuaweiCloudMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiCloudMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HuaweiCloudMachine{}, &HuaweiCloudMachineList{})
}
