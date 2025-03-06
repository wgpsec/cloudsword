package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
	"sync"
)

var ecsInstances []ecs.Instance

// ecs_list_instances

func ECSListInstances() {
	regions, err := describeRegions()
	if err == nil {
		ecsInstances = nil
		var (
			wg      sync.WaitGroup
			threads = 3
		)
		// 获取 ECS 资源
		taskCh := make(chan string, threads)
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				err := describeEcsInstances(taskCh, &wg)
				if err != nil {
					return
				}
			}()
		}
		for _, item := range regions {
			taskCh <- item
		}
		close(taskCh)
		wg.Wait()

		// 输出 ECS 资源
		if len(ecsInstances) == 0 {
			logger.Println.Warn("未找到 ECS 弹性计算实例，出现这种情况一般是由于没有实例或没有权限。\n")
		} else {
			logger.Println.Info("找到以下 ECS 弹性计算实例：")
			for n, i := range ecsInstances {
				fmt.Printf("\n=================== %d %v ===================\n", n+1, i.InstanceId)
				// 提取 ECS 的 IP 地址
				var (
					PublicIpAddress     string
					PublicIpAddressList []string
					PrivateIpAddress    string
				)
				if len(i.PublicIpAddress.IpAddress) > 0 {
					PublicIpAddressList = append(PublicIpAddressList, i.PublicIpAddress.IpAddress...)
				}
				if len(i.EipAddress.IpAddress) > 0 {
					PublicIpAddressList = append(PublicIpAddressList, i.EipAddress.IpAddress)
				}
				if len(PublicIpAddressList) > 0 {
					PublicIpAddress = strings.Join(PublicIpAddressList, ",")
				}

				if len(i.NetworkInterfaces.NetworkInterface[0].PrivateIpSets.PrivateIpSet) > 0 {
					PrivateIpAddress = i.NetworkInterfaces.NetworkInterface[0].PrivateIpSets.PrivateIpSet[0].PrivateIpAddress
				}

				if global.GetBasicOptionValue(global.Detail) == global.False {
					// 输出 ECS 基础信息
					fmt.Printf("\n实例 ID：%v\n", i.InstanceId)
					fmt.Printf("实例名称：%v\n", i.InstanceName)
					fmt.Printf("实例描述：%v\n", i.Description)
					fmt.Printf("实例区域：%v\n", i.RegionId)
					fmt.Printf("实例状态：%v\n", i.Status)
					fmt.Printf("操作系统类型：%v\n", i.OSType)
					fmt.Printf("私有 IP 地址：%v\n", PrivateIpAddress)
					fmt.Printf("公有 IP 地址：%v\n", PublicIpAddress)
				} else {
					// 输出 ECS 详细信息
					fmt.Printf("\n基础信息\n========\n")
					fmt.Printf("实例 ID：%v\n", i.InstanceId)
					fmt.Printf("实例名称：%v\n", i.InstanceName)
					fmt.Printf("实例描述：%v\n", i.Description)
					fmt.Printf("实例区域：%v\n", i.RegionId)
					fmt.Printf("元数据 HttpTokens: %v\n", i.MetadataOptions.HttpTokens)
					fmt.Printf("元数据 HttpEndpoint: %v\n", i.MetadataOptions.HttpEndpoint)

					ramRole, err := describeInstanceRamRole(i.RegionId, i.InstanceId)
					if err == nil {
						fmt.Printf("实例角色：%v\n", ramRole)
					}

					userData, err := describeUserData(i.RegionId, i.InstanceId)
					if err == nil {
						fmt.Printf("实例用户数据：%v\n", userData)
					}

					cloudAssistantStatus, err := describeCloudAssistantStatus(i.RegionId, i.InstanceId, i.OSType)
					if err == nil {
						fmt.Printf("实例云助手状态：%v\n", cloudAssistantStatus)
					}

					fmt.Printf("\n操作系统信息\n============\n")
					fmt.Printf("虚拟 CPU 数量：%v\n", i.Cpu)
					fmt.Printf("内存大小（单位：MiB）：%v\n", i.Memory)
					fmt.Printf("实例规格族：%v\n", i.InstanceTypeFamily)
					fmt.Printf("实例规格：%v\n", i.InstanceType)
					fmt.Printf("操作系统类型：%v\n", i.OSType)
					fmt.Printf("操作系统名称：%v\n", i.OSName)
					fmt.Printf("操作系统镜像 ID：%v\n", i.ImageId)

					fmt.Printf("\n网络信息\n========\n")
					fmt.Printf("网络类型：%v\n", i.InstanceNetworkType)
					fmt.Printf("VPC ID：%v\n", i.VpcAttributes.VpcId)
					fmt.Printf("安全组：%v\n", strings.Join(i.SecurityGroupIds.SecurityGroupId, ","))
					fmt.Printf("私有 IP 地址：%v\n", PrivateIpAddress)
					fmt.Printf("公有 IP 地址：%v\n", PublicIpAddress)
					fmt.Printf("公网入带宽最大值（单位：Mbit/s）：%v\n", i.InternetMaxBandwidthIn)
					fmt.Printf("公网出带宽最大值（单位：Mbit/s）：%v\n", i.InternetMaxBandwidthOut)

					fmt.Printf("\n其他信息========\n")
					fmt.Printf("删除保护：%v\n", i.DeletionProtection)
					jsonBytes, err := json.Marshal(i.Tags.Tag)
					if err != nil {
						logger.Println.Error(err.Error())
					}
					fmt.Printf("实例标签：%v\n", jsonBytes)
					fmt.Printf("计费类型：%v\n", i.InstanceChargeType)
					fmt.Printf("主机名称：%v\n", i.Hostname)
					fmt.Printf("资源组 ID：%v\n", i.ResourceGroupId)
					fmt.Printf("实例所属可用区：%v\n", i.ZoneId)
					fmt.Printf("实例序列号：%v\n", i.SerialNumber)

					creationTime, err := utils.GetUTC8TimeType1(i.CreationTime)
					if err == nil {
						fmt.Printf("创建时间：%v\n", creationTime)
					}

					expiredTime, err := utils.GetUTC8TimeType1(i.ExpiredTime)
					if err == nil {
						fmt.Printf("过期时间：%v\n", expiredTime)
					}

					fmt.Println(strings.Repeat("=", 60))
				}
			}
			fmt.Println()
		}
	}
}

func describeEcsInstances(ch <-chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var (
		err      error
		response *ecs.DescribeInstancesResponse
	)
	for region := range ch {
		logger.Println.Info(fmt.Sprintf("正在获取 %s 区域下的阿里云 ECS 资源信息", region))
		request := ecs.CreateDescribeInstancesRequest()
		for {
			client, err := ECSClient(region)
			if err == nil {
				response, err = client.DescribeInstances(request)
				if err != nil {
					break
				}
				if len(response.Instances.Instance) > 0 {
					logger.Println.Warn(fmt.Sprintf("在 %s 区域下获取到 %d 条 ECS 资源", region, len(response.Instances.Instance)))
					ecsInstances = append(ecsInstances, response.Instances.Instance...)
				}
				if response.NextToken == "" {
					break
				}
				request.NextToken = response.NextToken
			}
		}
	}
	return err
}
