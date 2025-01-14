package securitygroup

import (
	"k8s.io/klog/v2"

	"github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pkg/scope"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/region"
)

type Service struct {
	scope     *scope.ClusterScope
	vpcClient *vpc.VpcClient
}

func NewService(scope *scope.ClusterScope) (*Service, error) {
	region, err := region.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}

	vpcHCHttpCli, err := vpc.VpcClientBuilder().
		WithRegion(region).
		WithCredential(scope.Credentials).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create VPC client: %v", err)
		return nil, err
	}
	vpcCli := vpc.NewVpcClient(vpcHCHttpCli)

	return &Service{
		scope:     scope,
		vpcClient: vpcCli,
	}, nil
}
