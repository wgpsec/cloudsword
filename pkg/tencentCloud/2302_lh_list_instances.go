package tencentCloud

import (
	"encoding/json"
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

// lh_list_instances

func LHListInstances() {
	lhInstances, err := GetLHResource()
	if err == nil {
		// 输出 LH 资源
		if len(lhInstances) == 0 {
			logger.Println.Warn("未找到 LH 弹性计算实例，出现这种情况一般是由于没有实例或没有权限。\n")
		} else {
			logger.Println.Info("找到以下 LH 弹性计算实例：")
			for n, i := range lhInstances {
				fmt.Printf("\n=================== %d %v ===================\n", n+1, *i.InstanceId)
				// 提取 LH 的 IP 地址
				var (
					PublicIpAddressList []string
					PrivateIpAddress    string
					PublicIpAddress     string
				)
				if len(i.PublicAddresses) > 0 {
					for _, v := range i.PublicAddresses {
						PublicIpAddressList = append(PublicIpAddressList, *v)
					}
					PublicIpAddress = strings.Join(PublicIpAddressList, ",")
				}
				if len(i.PrivateAddresses) > 0 {
					PrivateIpAddress = *i.PrivateAddresses[0]
				}

				if global.GetBasicOptionValue(global.Detail) == global.False {
					// 输出 ECS 基础信息
					utils.PrintfNotNilString("\n实例 ID：", i.InstanceId)
					utils.PrintfNotNilString("实例名称：", i.InstanceName)
					utils.PrintfNotNilString("实例区域：", i.Zone)
					utils.PrintfNotNilString("实例状态：", i.InstanceState)
					utils.PrintfNotNilString("操作系统类型：", i.OsName)
					fmt.Printf("私有 IP 地址：%v\n", PrivateIpAddress)
					fmt.Printf("公有 IP 地址：%v\n", PublicIpAddress)
				} else {
					// 输出 ECS 详细信息
					fmt.Printf("\n基础信息\n========\n")
					utils.PrintfNotNilString("实例 ID：", i.InstanceId)
					utils.PrintfNotNilString("实例名称：", i.InstanceName)
					utils.PrintfNotNilString("实例区域：", i.Zone)
					utils.PrintfNotNilString("实例状态：", i.InstanceState)

					fmt.Printf("\n操作系统信息\n============\n")
					utils.PrintfNotNilInt64("虚拟 CPU 数量：", i.CPU)
					utils.PrintfNotNilInt64("内存大小（单位：GB）：", i.Memory)
					utils.PrintfNotNilString("操作系统类型：", i.OsName)
					utils.PrintfNotNilString("操作系统平台：", i.Platform)
					utils.PrintfNotNilString("操作平台类型：", i.PlatformType)

					fmt.Printf("\n磁盘信息\n============\n")
					utils.PrintfNotNilString("磁盘 ID：", i.SystemDisk.DiskId)
					utils.PrintfNotNilInt64("磁盘大小（单位：GB）：", i.SystemDisk.DiskSize)
					utils.PrintfNotNilString("磁盘类型：", i.SystemDisk.DiskType)

					fmt.Printf("\n网络信息\n========\n")
					fmt.Printf("私有 IP 地址：%v\n", PrivateIpAddress)
					fmt.Printf("公有 IP 地址：%v\n", PublicIpAddress)
					utils.PrintfNotNilString("网络计费类型：", i.InternetAccessible.InternetChargeType)
					utils.PrintfNotNilInt64("公网出带宽最大值（单位：Mbps）：", i.InternetAccessible.InternetMaxBandwidthOut)

					fmt.Printf("\n其他信息========\n")
					jsonBytes, err := json.Marshal(i.Tags)
					if err == nil {
						fmt.Printf("实例标签：%s\n", jsonBytes)
					}

					utils.PrintfNotNilString("计费类型：", i.InstanceChargeType)
					utils.PrintfNotNilString("实例全局唯一ID：", i.Uuid)

					if i.CreatedTime != nil {
						creationTime, err := utils.GetUTC8TimeType2(*i.CreatedTime)
						if err == nil {
							fmt.Printf("创建时间：%v\n", creationTime)
						}
					}

					if i.ExpiredTime != nil {
						expiredTime, err := utils.GetUTC8TimeType2(*i.ExpiredTime)
						if err == nil {
							fmt.Printf("过期时间：%v\n", expiredTime)
						}
					}
				}
			}
			fmt.Println()
		}
	}
}
