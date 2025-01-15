package elb

import (
	"k8s.io/klog/v2"

	"github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pkg/scope"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	eipregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/region"
	elbregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v2/region"
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
)

type Service struct {
	scope     *scope.ClusterScope
	elbClient *elb.ElbClient
	eipClient *eip.EipClient
}

func NewService(scope *scope.ClusterScope) (*Service, error) {
	elbreg, err := elbregion.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}

	elbHCHttpCli, err := elb.ElbClientBuilder().
		WithRegion(elbreg).
		WithCredential(scope.Credentials).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create ELB client: %v", err)
		return nil, err
	}

	elbCli := elb.NewElbClient(elbHCHttpCli)

	eipReg, err := eipregion.SafeValueOf(scope.Region())
	if err != nil {
		klog.Errorf("Failed to get region: %v", err)
		return nil, err
	}

	eipHCHttpCli, err := eip.EipClientBuilder().
		WithRegion(eipReg).
		WithCredential(scope.Credentials).
		SafeBuild()
	if err != nil {
		klog.Errorf("Failed to create EIP client: %v", err)
		return nil, err
	}
	eipCli := eip.NewEipClient(eipHCHttpCli)

	return &Service{
		elbClient: elbCli,
		eipClient: eipCli,
		scope:     scope,
	}, nil
}
