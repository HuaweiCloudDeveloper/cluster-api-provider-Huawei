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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HuaweiCloudMachineTemplateSpec defines the desired state of HuaweiCloudMachineTemplate.
type HuaweiCloudMachineTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HuaweiCloudMachineTemplate. Edit huaweicloudmachinetemplate_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// HuaweiCloudMachineTemplateStatus defines the observed state of HuaweiCloudMachineTemplate.
type HuaweiCloudMachineTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// HuaweiCloudMachineTemplate is the Schema for the huaweicloudmachinetemplates API.
type HuaweiCloudMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HuaweiCloudMachineTemplateSpec   `json:"spec,omitempty"`
	Status HuaweiCloudMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HuaweiCloudMachineTemplateList contains a list of HuaweiCloudMachineTemplate.
type HuaweiCloudMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiCloudMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HuaweiCloudMachineTemplate{}, &HuaweiCloudMachineTemplateList{})
}
