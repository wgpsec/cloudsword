package huaweiCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

func OBSListBuckets() {
	buckets, err := listOBSBuckets()
	if err == nil {
		if len(buckets) == 0 {
			logger.Println.Info("没有找到存储桶。\n")
		} else {
			logger.Println.Info("找到以下存储桶：")
			for _, b := range buckets {
				if global.GetBasicOptionValue(global.Detail) == global.False {
					fmt.Println(b.Name)
				} else {
					fmt.Printf("\n名称：%v\n", b.Name)
					fmt.Printf("区域：%v\n", b.Location)
					fmt.Printf("访问地址：https://%v.obs.%v.myhuaweicloud.com\n", b.Name, b.Location)
					creationDate := utils.FormatTime(&b.CreationDate)
					if err == nil {
						fmt.Printf("创建时间：%v\n", creationDate)
					}
					fmt.Println(strings.Repeat("=", 60))
				}
			}
			fmt.Println()
		}
	}
}
