package tencentCloud

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"sync"
)

var CVMInstances []*cvm.Instance

// 获取 CVM 区域

func CVMDescribeRegions() ([]string, error) {
	var regions []string
	CVMInstances = nil
	region := global.GetBasicOptionValue(global.Region)
	if region == global.All {
		client, err := CVMClient(global.TencentCloudDefaultRegion)
		if err == nil {
			cvmRequest := cvm.NewDescribeRegionsRequest()
			cvmRequest.SetScheme("https")
			cvmResponse, err := client.DescribeRegions(cvmRequest)
			if err != nil {
				logger.Println.Error("获取 CVM 区域时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				return regions, err
			} else {
				for _, region := range cvmResponse.Response.RegionSet {
					regions = append(regions, *region.Region)
				}
				return regions, err
			}
		} else {
			return regions, err
		}
	} else {
		regions = append(regions, region)
		return regions, nil
	}
}

// 获取 CVM 弹性计算实例

func GetCVMResource() ([]*cvm.Instance, error) {
	regions, err := CVMDescribeRegions()
	if err == nil {
		var (
			threads = 3
			wg      sync.WaitGroup
		)
		taskCh := make(chan string, threads)
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				err = CVMDescribeInstances(taskCh, &wg)
			}()
		}
		for _, item := range regions {
			taskCh <- item
		}
		close(taskCh)
		wg.Wait()
		return CVMInstances, err
	} else {
		return nil, err
	}
}

func CVMDescribeInstances(ch <-chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var (
		err      error
		client   *cvm.Client
		response *cvm.DescribeInstancesResponse
		limit    int64 = 100
		offset   int64 = 0
	)
	for region := range ch {
		logger.Println.Info(fmt.Sprintf("正在获取 %s 区域下的腾讯云 CVM 资源信息", region))
		client, err = CVMClient(region)
		if err == nil {
			for {
				request := cvm.NewDescribeInstancesRequest()
				request.Limit = common.Int64Ptr(limit)
				request.SetScheme("https")
				response, err = client.DescribeInstances(request)
				if err != nil {
					logger.Println.Error(fmt.Sprintf("获取 %v 区域下的 CVM 弹性计算实例报错，详细信息如下：", region))
					logger.Println.Error(err.Error())
					break
				} else {
					if len(response.Response.InstanceSet) > 0 {
						logger.Println.Warn(fmt.Sprintf("在 %s 区域下获取到 %d 条 CVM 资源", region, len(response.Response.InstanceSet)))
						CVMInstances = append(CVMInstances, response.Response.InstanceSet...)
					}
					if len(response.Response.InstanceSet) == int(limit) {
						offset = offset + limit
						request.Offset = common.Int64Ptr(offset)
					} else {
						break
					}
				}
			}
		}
	}
	return err
}
