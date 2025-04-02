package qiniuCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

func KodoListBuckets() {
	buckets, err := listKodoBuckets()
	if err == nil {
		if len(buckets) == 0 {
			logger.Println.Info("没有找到存储桶。\n")
		} else {
			logger.Println.Info("找到以下存储桶：")
			for _, b := range buckets {
				if global.GetBasicOptionValue(global.Detail) == global.False {
					fmt.Println(b)
				} else {
					fmt.Printf("\n名称：%v\n", b)
					fmt.Println(strings.Repeat("=", 60))
				}
			}
			fmt.Println()
		}
	}
}
