package tencentCloud

import (
	"encoding/json"
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

// cvm_list_instances

func CVMListInstances() {
	cvmInstances, err := GetCVMResource()
	if err == nil {
		// 输出 CVM 资源
		if len(cvmInstances) == 0 {
			logger.Println.Warn("未找到 CVM 弹性计算实例，出现这种情况一般是由于没有实例或没有权限。\n")
		} else {
			logger.Println.Info("找到以下 CVM 弹性计算实例：")
			for n, i := range cvmInstances {
				fmt.Printf("\n=================== %d %v ===================\n", n+1, *i.InstanceId)
				// 提取 CVM 的 IP 地址
				var (
					PublicIpAddressList []string
					PrivateIpAddress    string
					PublicIpAddress     string
				)
				if len(i.PublicIpAddresses) > 0 {
					for _, v := range i.PublicIpAddresses {
						PublicIpAddressList = append(PublicIpAddressList, *v)
					}
					PublicIpAddress = strings.Join(PublicIpAddressList, ",")
				}
				if len(i.PrivateIpAddresses) > 0 {
					PrivateIpAddress = *i.PrivateIpAddresses[0]
				}

				if global.GetBasicOptionValue(global.Detail) == global.False {
					// 输出 CVM 基础信息
					utils.PrintfNotNilString("\n实例 ID：", i.InstanceId)
					utils.PrintfNotNilString("实例名称：", i.InstanceName)
					utils.PrintfNotNilString("实例区域：", i.Placement.Zone)
					utils.PrintfNotNilString("实例状态：", i.InstanceState)
					utils.PrintfNotNilString("操作系统类型：", i.OsName)
					fmt.Printf("私有 IP 地址：%v\n", PrivateIpAddress)
					fmt.Printf("公有 IP 地址：%v\n", PublicIpAddress)
					utils.PrintfNotNilInt64("默认登录端口：", i.DefaultLoginPort)
					utils.PrintfNotNilString("默认登录用户：", i.DefaultLoginUser)
					utils.PrintfNotNilString("登录密码：", i.LoginSettings.Password)
				} else {
					// 输出 CVM 详细信息
					fmt.Printf("\n基础信息\n========\n")
					utils.PrintfNotNilString("实例 ID：", i.InstanceId)
					utils.PrintfNotNilString("实例名称：", i.InstanceName)
					utils.PrintfNotNilString("实例区域：", i.Placement.Zone)
					utils.PrintfNotNilString("实例状态：", i.InstanceState)
					utils.PrintfNotNilString("实例角色：", i.CamRoleName)
					utils.PrintfNotNilInt64("默认登录端口：", i.DefaultLoginPort)
					utils.PrintfNotNilString("默认登录用户：", i.DefaultLoginUser)
					utils.PrintfNotNilString("登录密码：", i.LoginSettings.Password)

					fmt.Printf("\n操作系统信息\n============\n")
					utils.PrintfNotNilInt64("虚拟 CPU 数量：", i.CPU)
					utils.PrintfNotNilInt64("内存大小（单位：GB）：", i.Memory)
					if i.GPUInfo != nil {
						utils.PrintfNotNilFloat64("GPU 数量：", i.GPUInfo.GPUCount)
						utils.PrintfNotNilString("GPU 类型：", i.GPUInfo.GPUType)
					}
					utils.PrintfNotNilString("实例规格：", i.InstanceType)
					utils.PrintfNotNilString("操作系统类型：", i.OsName)
					utils.PrintfNotNilString("操作系统镜像 ID：", i.ImageId)

					fmt.Printf("\n磁盘信息\n============\n")
					utils.PrintfNotNilString("磁盘 ID：", i.SystemDisk.DiskId)
					utils.PrintfNotNilString("磁盘名称：", i.SystemDisk.DiskName)
					utils.PrintfNotNilInt64("磁盘大小（单位：GB）：", i.SystemDisk.DiskSize)
					utils.PrintfNotNilString("磁盘类型：", i.SystemDisk.DiskType)
					fmt.Println()
					for _, disk := range i.DataDisks {
						utils.PrintfNotNilString("磁盘 ID：", disk.DiskId)
						utils.PrintfNotNilString("磁盘名称：", disk.DiskName)
						utils.PrintfNotNilInt64("磁盘大小（单位：GB）：", disk.DiskSize)
						utils.PrintfNotNilString("磁盘类型：", disk.DiskType)
						utils.PrintfNotNilString("快照 ID：", disk.SnapshotId)
						fmt.Println()
					}

					fmt.Printf("\n网络信息\n========\n")
					utils.PrintfNotNilString("VPC ID：", i.VirtualPrivateCloud.VpcId)
					utils.PrintfNotNilString("子网 ID：", i.VirtualPrivateCloud.SubnetId)
					var securityGroupIds []string
					for _, v := range i.SecurityGroupIds {
						securityGroupIds = append(securityGroupIds, *v)
					}
					fmt.Printf("安全组：%v\n", strings.Join(securityGroupIds, ","))
					fmt.Printf("私有 IP 地址：%v\n", PrivateIpAddress)
					fmt.Printf("公有 IP 地址：%v\n", PublicIpAddress)
					utils.PrintfNotNilString("带宽包 ID：", i.InternetAccessible.BandwidthPackageId)
					utils.PrintfNotNilString("网络计费类型：", i.InternetAccessible.InternetChargeType)
					utils.PrintfNotNilInt64("公网出带宽最大值（单位：Mbps）：", i.InternetAccessible.InternetMaxBandwidthOut)

					fmt.Printf("\n其他信息========\n")
					jsonBytes, err := json.Marshal(i.Tags)
					if err == nil {
						fmt.Printf("实例标签：%s\n", jsonBytes)
					}
					utils.PrintfNotNilString("计费类型：", i.InstanceChargeType)
					utils.PrintfNotNilString("许可类型：", i.LicenseType)
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
