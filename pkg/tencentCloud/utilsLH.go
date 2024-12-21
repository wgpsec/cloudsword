package tencentCloud

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	lh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"sync"
)

var LHInstances []*lh.Instance

// 获取 LH 区域

func LHDescribeRegions() ([]string, error) {
	var regions []string
	LHInstances = nil
	region := global.GetBasicOptionValue(global.Region)
	if region == global.All {
		client, err := LHClient(global.TencentCloudDefaultRegion)
		if err == nil {
			lhRequest := lh.NewDescribeRegionsRequest()
			lhRequest.SetScheme("https")
			lhResponse, err := client.DescribeRegions(lhRequest)
			if err != nil {
				logger.Println.Error("获取 LH 区域时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				return regions, err
			} else {
				for _, region := range lhResponse.Response.RegionSet {
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

// 获取 LH 弹性计算实例

func GetLHResource() ([]*lh.Instance, error) {
	regions, err := LHDescribeRegions()
	if err == nil {
		var (
			threads = 3
			wg      sync.WaitGroup
		)
		taskCh := make(chan string, threads)
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				err = LHDescribeInstances(taskCh, &wg)
			}()
		}
		for _, item := range regions {
			taskCh <- item
		}
		close(taskCh)
		wg.Wait()
		return LHInstances, err
	} else {
		return nil, err
	}
}

func LHDescribeInstances(ch <-chan string, wg *sync.WaitGroup) error {
	defer wg.Done()
	var (
		err      error
		client   *lh.Client
		response *lh.DescribeInstancesResponse
		limit    int64 = 100
		offset   int64 = 0
	)
	for region := range ch {
		logger.Println.Info(fmt.Sprintf("正在获取 %s 区域下的腾讯云 LH 资源信息", region))
		client, err = LHClient(region)
		if err == nil {
			for {
				request := lh.NewDescribeInstancesRequest()
				request.Limit = common.Int64Ptr(limit)
				request.SetScheme("https")
				response, err = client.DescribeInstances(request)
				if err != nil {
					logger.Println.Error(fmt.Sprintf("获取 %v 区域下的 LH 弹性计算实例报错，详细信息如下：", region))
					logger.Println.Error(err.Error())
					break
				} else {
					if len(response.Response.InstanceSet) > 0 {
						logger.Println.Warn(fmt.Sprintf("在 %s 区域下获取到 %d 条 LH 资源", region, len(response.Response.InstanceSet)))
						LHInstances = append(LHInstances, response.Response.InstanceSet...)
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
