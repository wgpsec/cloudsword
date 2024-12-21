package aliyun

import (
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

// oss_list_buckets

func OSSListBuckets() {
	bucket := listBuckets()
	if len(bucket) == 0 {
		logger.Println.Info("没有找到存储桶。\n")
	} else {
		logger.Println.Info("找到以下存储桶：")
		for _, b := range bucket {
			if global.GetBasicOptionValue(global.Detail) == global.False {
				fmt.Println(oss.ToString(b.Name))
			} else {
				bucketsStat, err := getBucketStat(oss.ToString(b.Name), oss.ToString(b.Region))
				fmt.Printf("\n基础信息\n========\n")
				fmt.Printf("名称：%v\n", oss.ToString(b.Name))
				fmt.Printf("区域：%v\n", oss.ToString(b.Region))
				if err == nil {
					fmt.Printf("存储对象数量：%v\n", bucketsStat.ObjectCount)
					fmt.Printf("存储容量：%v\n", utils.FormatBytes(bucketsStat.Storage))
				}
				bucketTags, err := getBucketTags(oss.ToString(b.Name), oss.ToString(b.Region))
				if err == nil {
					fmt.Printf("存储桶标签：%v\n", bucketTags)
				}

				fmt.Printf("访问地址：https://%v.oss-%v.aliyuncs.com\n", oss.ToString(b.Name), oss.ToString(b.Region))

				fmt.Printf("\n更多信息\n========\n")
				fmt.Printf("存储类型：%v\n", oss.ToString(b.StorageClass))
				fmt.Printf("数据中心位置：%v\n", oss.ToString(b.Location))
				fmt.Printf("创建时间：%v\n", utils.FormatTime(b.CreationDate))
				fmt.Printf("公共端点：%v\n", oss.ToString(b.ExtranetEndpoint))
				fmt.Printf("内部端点：%v\n", oss.ToString(b.IntranetEndpoint))
				fmt.Printf("资源组 ID：%v\n", oss.ToString(b.ResourceGroupId))
				fmt.Println(strings.Repeat("=", 60))
			}
		}
		fmt.Println()
	}
}
