/*
Copyright 2024.

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

// NetworkSpect encapsulates the configuration options for HuaweiCloud network.
type NetworkSpec struct {
	// VPC configuration.
	// +optional
	VPC VPCSpec `json:"vpc,omitempty"`
}

type VPCSpec struct {
	// Id is the unique identifier of the VPC. It is a UUID.
	Id string `json:"id"`

	// Name is the name of the VPC. It must be 0-64 characters long and support numbers, letters, Chinese characters, _(underscore), -(hyphen), and .(dot).
	Name string `json:"name"`

	// Cidr is the CIDR of the VPC.
	Cidr string `json:"cidr"`
}
