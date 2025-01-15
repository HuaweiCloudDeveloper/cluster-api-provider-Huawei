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

	// Subnets configuration.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`
}

type VPCSpec struct {
	// Id is the unique identifier of the VPC. It is a UUID.
	Id string `json:"id"`

	// Name is the name of the VPC. It must be 0-64 characters long and support numbers, letters, Chinese characters, _(underscore), -(hyphen), and .(dot).
	Name string `json:"name"`

	// Cidr is the CIDR of the VPC.
	Cidr string `json:"cidr"`
}

type Subnets []SubnetSpec

// Subnet
type SubnetSpec struct {
	// ID defines a unique identifier to reference this resource.
	Id string `json:"id"`

	// Name is the name of the subnet. It must be 1-64 characters long and support numbers, letters, Chinese characters, _(underscore), -(hyphen), and .(dot).
	Name string `json:"name"`

	// CIDR is the CIDR of the subnet. It must be in CIDR format. The mask length cannot be greater than 28.
	Cidr string `json:"cidr"`

	// GatewayIp is the gateway of the subnet. It must be an IP address in the subnet segment.
	GatewayIp string `json:"gateway_ip"`

	// VPCId is the identifier of the VPC where the subnet is located.
	VpcId string `json:"vpc_id"`

	// NeutronNetworkId is the identifier of the network (OpenStack Neutron interface).
	NeutronNetworkId string `json:"neutron_network_id"`

	// NeutronSubnetId is the identifier of the subnet (OpenStack Neutron interface).
	NeutronSubnetId string `json:"neutron_subnet_id"`
}
