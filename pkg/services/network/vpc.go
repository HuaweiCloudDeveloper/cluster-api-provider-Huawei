package network

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

func (s *Service) reconcileVPC() error {
	// Check if VPC exists, if not create it
	request := &model.ListVpcsRequest{}
	response, err := s.vpcClient.ListVpcs(request)
	if err != nil {
		return errors.Wrap(err, "failed to list VPCs")
	}

	var vpc *model.Vpc
	if len(*response.Vpcs) == 0 {

		createRequest := &model.CreateVpcRequest{}
		cidrVpc := "192.168.0.0/16"
		nameVpc := "k8s-vpc"
		vpcbody := &model.CreateVpcOption{
			Cidr: &cidrVpc,
			Name: &nameVpc,
		}
		createRequest.Body = &model.CreateVpcRequestBody{
			Vpc: vpcbody,
		}

		createRes, err := s.vpcClient.CreateVpc(createRequest)
		if err != nil {
			return errors.Wrap(err, "failed to create VPC")
		}
		vpc = createRes.Vpc
		klog.Infof("Created VPC %s", vpc.Id)
	} else {
		vpc = &(*response.Vpcs)[0]
		klog.Infof("VPC %s already exists", vpc.Id)
	}

	s.scope.VPC().Id = vpc.Id
	s.scope.VPC().Name = vpc.Name
	s.scope.VPC().Cidr = vpc.Cidr

	if err := s.scope.PatchObject(); err != nil {
		return errors.Wrap(err, "failed to patch HuaweiCloudCluster with PVC details")
	}

	return nil
}

func (s *Service) deleteVPC() error {
	if s.scope.VPC().Id == "" {
		return errors.New("VPC ID is empty")
	}

	deleteRequest := &model.DeleteVpcRequest{
		VpcId: s.scope.VPC().Id,
	}
	_, err := s.vpcClient.DeleteVpc(deleteRequest)
	if err != nil {
		return errors.Wrapf(err, "failed to delete VPC %s", s.scope.VPC().Id)
	}

	klog.Infof("Deleted VPC %s", s.scope.VPC().Id)
	return nil
}
