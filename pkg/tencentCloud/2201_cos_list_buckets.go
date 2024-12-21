package tencentCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

// cos_list_buckets

func COSListBuckets() {
	buckets := listCOSBuckets()
	if len(buckets) == 0 {
		logger.Println.Info("没有找到存储桶。\n")
	} else {
		logger.Println.Info("找到以下存储桶：")
		for _, b := range buckets {
			if global.GetBasicOptionValue(global.Detail) == global.False {
				fmt.Println(b.Name)
			} else {
				fmt.Printf("\n名称：%v\n", b.Name)
				fmt.Printf("区域：%v\n", b.Region)

				bucketTags, err := getBucketTagging(b.Name, b.Region)
				if err == nil {
					fmt.Printf("存储桶标签：%v\n", bucketTags)
				}

				fmt.Printf("访问地址：https://%v.cos.%v.myqcloud.com\n", b.Name, b.Region)
				creationDate, err := utils.GetUTC8TimeType2(b.CreationDate)
				if err == nil {
					fmt.Printf("创建时间：%v\n", creationDate)
				}
				fmt.Println(strings.Repeat("=", 60))
			}
		}
		fmt.Println()
	}
}
