package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

func describeCloudAssistantStatus(region, instanceId, OSType string) (string, error) {
	request := ecs.CreateDescribeCloudAssistantStatusRequest()
	request.RegionId = region
	request.OSType = OSType
	request.InstanceId = &[]string{instanceId}
	client, err := ECSClient(region)
	if err == nil {
		response, err := client.DescribeCloudAssistantStatus(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("获取 %v 实例云助手状态时报错，详细信息如下：", instanceId))
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			if len(response.InstanceCloudAssistantStatusSet.InstanceCloudAssistantStatus) > 0 {
				cloudAssistantStatus := response.InstanceCloudAssistantStatusSet.InstanceCloudAssistantStatus[0]
				if cloudAssistantStatus.CloudAssistantStatus == "true" {
					return "已安装", err
				} else {
					return "未安装", err
				}
			} else {
				return "未安装", err
			}
		}
	} else {
		return global.NULL, err
	}
}

func describeInstanceRamRole(region, instanceId string) (string, error) {
	var roles []string
	request := ecs.CreateDescribeInstanceRamRoleRequest()
	request.RegionId = region
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", instanceId)
	client, err := ECSClient(region)
	if err == nil {
		response, err := client.DescribeInstanceRamRole(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("获取 %v ECS 实例 RAM 角色时报错，详细信息如下：", instanceId))
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			for _, role := range response.InstanceRamRoleSets.InstanceRamRoleSet {
				roles = append(roles, role.RamRoleName)
			}
			return strings.Join(roles, ","), err
		}
	} else {
		return global.NULL, err
	}
}

func describeUserData(region, instanceId string) (string, error) {
	request := ecs.CreateDescribeUserDataRequest()
	request.RegionId = region
	request.InstanceId = instanceId
	client, err := ECSClient(region)
	if err == nil {
		response, err := client.DescribeUserData(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("获取 %v 实例用户数据时报错，详细信息如下：", instanceId))
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			return response.UserData, err
		}
	} else {
		return global.NULL, err
	}
}

// 获取 ECS 区域

func describeRegions() ([]string, error) {
	var regions []string
	globalRegion := global.GetBasicOptionValue(global.Region)
	if globalRegion == global.All {
		client, err := ECSClient(global.AliyunDefaultRegion)
		if err == nil {
			ecsRegions, err := client.DescribeRegions(ecs.CreateDescribeRegionsRequest())
			if err != nil {
				logger.Println.Error("获取 ECS 弹性计算实例区域时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				return regions, err
			} else {
				for _, region := range ecsRegions.Regions.Region {
					regions = append(regions, region.RegionId)
				}
				return regions, err
			}
		} else {
			return regions, err
		}
	} else {
		regions = append(regions, globalRegion)
		return regions, nil
	}
}
