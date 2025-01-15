package elb

import (
	"fmt"
	"strings"

	infrav1alpha1 "github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/api/v1alpha1"
	eipmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
	elbmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

type HuaweiElbPool struct {
	Id   string `json:"id,omitempty"`
	Port int32  `json:"port,omitempty"`
}

func getLoadBalancerChargeMode(chargeMode string) elbmodel.CreateLoadBalancerBandwidthOptionChargeMode {
	var chargeModeEnum elbmodel.CreateLoadBalancerBandwidthOptionChargeMode
	switch chargeMode {
	case "traffic":
		chargeModeEnum = elbmodel.GetCreateLoadBalancerBandwidthOptionChargeModeEnum().TRAFFIC
	case "bandwidth":
		chargeModeEnum = elbmodel.GetCreateLoadBalancerBandwidthOptionChargeModeEnum().BANDWIDTH
	default:
		chargeModeEnum = elbmodel.GetCreateLoadBalancerBandwidthOptionChargeModeEnum().TRAFFIC
	}
	return chargeModeEnum
}

func getLoadBalancerShareType(shareType string) elbmodel.CreateLoadBalancerBandwidthOptionShareType {
	var shareTypeEnum elbmodel.CreateLoadBalancerBandwidthOptionShareType
	switch shareType {
	case "per":
		shareTypeEnum = elbmodel.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().PER
	case "whole":
		shareTypeEnum = elbmodel.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().WHOLE
	default:
		shareTypeEnum = elbmodel.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().PER
	}
	return shareTypeEnum
}

// ReconcileLoadbalancers reconciles the load balancers for the given cluster.
func (s *Service) ReconcileLoadbalancers() error {
	klog.Info("Reconciling load balancers")

	lbName := fmt.Sprintf("%s-elb", s.scope.ClusterName())
	lb, err := s.getLoadBalancerByName(lbName)
	if err != nil {
		return errors.Wrapf(err, "failed to get load balancer %s", lbName)
	}

	if lb == nil {
		klog.Info("Creating new load balancer", "name", lbName)
		if err := s.createLoadBalancer(lbName); err != nil {
			return errors.Wrapf(err, "failed to create load balancer %s", lbName)
		}
		// Re-fetch the load balancer after creation
		lb, err = s.getLoadBalancerByName(lbName)
		if err != nil {
			return errors.Wrapf(err, "failed to get load balancer %s after creation", lbName)
		}
	} else {
		klog.Info("Load balancer already exists", "name", lbName)
	}

	if lb != nil {
		s.scope.HCCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
			Host: lb.Publicips[0].PublicipAddress,
			Port: 6443,
		}
	}

	conditions.MarkTrue(s.scope.InfraCluster(), infrav1alpha1.LoadBalancerReadyCondition)
	return nil
}

// DeleteLoadbalancers deletes the load balancers for the given cluster.
func (s *Service) DeleteLoadbalancers() error {
	klog.Info("Deleting load balancers")

	lbName := fmt.Sprintf("%s-elb", s.scope.ClusterName())
	lb, err := s.getLoadBalancerByName(lbName)
	if err != nil {
		return errors.Wrapf(err, "failed to get load balancer %s", lbName)
	}

	if lb != nil {
		klog.Info("Deleting load balancer", "name", lbName)
		if err := s.deleteLoadBalancer(lb.Id); err != nil {
			conditions.MarkFalse(
				s.scope.InfraCluster(),
				infrav1alpha1.LoadBalancerReadyCondition,
				clusterv1.DeletingReason,
				clusterv1.ConditionSeverityWarning,
				"failed to delete load balancer")
			return errors.Wrapf(err, "failed to delete load balancer %s", lbName)
		}

		// delete related elastic ip
		for _, publicIp := range lb.Publicips {
			delPubIpReq := &eipmodel.DeletePublicipRequest{
				PublicipId: publicIp.PublicipId,
			}
			delPubIpRes, err := s.eipClient.DeletePublicip(delPubIpReq)
			if err != nil {
				conditions.MarkFalse(
					s.scope.InfraCluster(),
					infrav1alpha1.LoadBalancerReadyCondition,
					clusterv1.DeletingReason,
					clusterv1.ConditionSeverityWarning,
					"failed to delete public ip")
				return errors.Wrapf(err, "failed to delete public ip %s", publicIp.PublicipId)
			}
			klog.Infof("Delete public ip response: %v", delPubIpRes)
		}
	}

	conditions.MarkFalse(
		s.scope.InfraCluster(),
		infrav1alpha1.LoadBalancerReadyCondition,
		clusterv1.DeletedReason,
		clusterv1.ConditionSeverityInfo,
		"")
	return nil
}

func (s *Service) getLoadBalancerByName(name string) (*elbmodel.LoadBalancer, error) {
	names := []string{name}
	request := &elbmodel.ListLoadBalancersRequest{
		Name: &names,
	}

	response, err := s.elbClient.ListLoadBalancers(request)
	if err != nil {
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to list load balancers with name %s", name)
	}

	if response == nil || len(*response.Loadbalancers) == 0 {
		return nil, nil
	}

	return &(*response.Loadbalancers)[0], nil
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "APIGW.0101")
}

func (s *Service) getAvailabilityZones() ([]string, error) {
	availabilityZones := make([]string, 0)
	request := &elbmodel.ListAvailabilityZonesRequest{}
	response, err := s.elbClient.ListAvailabilityZones(request)
	if err != nil {
		return nil, err
	}
	for _, zones := range *response.AvailabilityZones {
		for _, zone := range zones {
			availabilityZones = append(availabilityZones, zone.Code)
			if len(availabilityZones) == 2 {
				break
			}
		}
	}
	return availabilityZones, nil
}

func (s *Service) createLoadBalancer(name string) error {
	request := &elbmodel.CreateLoadBalancerRequest{}
	nameBandwidth := "eip-caph"
	chargeMode := getLoadBalancerChargeMode("traffic")
	shareType := getLoadBalancerShareType("per")
	var bwSize int32 = 100
	bandwidthPublicip := &elbmodel.CreateLoadBalancerBandwidthOption{
		Name:       &nameBandwidth,
		Size:       &bwSize,
		ChargeMode: &chargeMode,
		ShareType:  &shareType,
	}
	publicipLoadbalancer := &elbmodel.CreateLoadBalancerPublicIpOption{
		NetworkType: "5_bgp",
		Bandwidth:   bandwidthPublicip,
	}
	zones, err := s.getAvailabilityZones()
	if err != nil {
		return err
	}
	loadbalancerbody := &elbmodel.CreateLoadBalancerOption{
		Name:                 &name,
		VipSubnetCidrId:      &s.scope.Subnets()[0].NeutronSubnetId,
		VpcId:                &s.scope.VPC().Id,
		AvailabilityZoneList: zones,
		Publicip:             publicipLoadbalancer,
	}
	request.Body = &elbmodel.CreateLoadBalancerRequestBody{
		Loadbalancer: loadbalancerbody,
	}
	klog.Infof("Create load balancer request: %v", request)
	response, err := s.elbClient.CreateLoadBalancer(request)
	if err != nil {
		return err
	}
	klog.Infof("Create load balancer response: %v", response)
	return nil
}

func (s *Service) deleteLoadBalancer(id string) error {
	request := &elbmodel.DeleteLoadBalancerRequest{
		LoadbalancerId: id,
	}

	_, err := s.elbClient.DeleteLoadBalancer(request)
	return err
}
