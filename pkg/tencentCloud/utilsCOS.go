package tencentCloud

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

// 列存储桶

func listCOSBuckets() []cos.Bucket {
	var (
		NextMarker string
		buckets    []cos.Bucket
	)
	for {
		response, _, err := COSClient(global.NULL, global.NULL, NextMarker).Service.Get(context.Background())
		if err != nil {
			logger.Println.Error("列出存储桶时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			fmt.Println()
			break
		} else {
			for _, bucket := range response.Buckets {
				buckets = append(buckets, bucket)
			}
			if response.IsTruncated {
				NextMarker = response.NextMarker
			} else {
				break
			}
		}
	}
	return buckets
}

// 获取存储桶标签

func getBucketTagging(bucketName, region string) (string, error) {
	response, _, err := COSClient(bucketName, region, global.NULL).Bucket.GetTagging(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "The TagSet does not exist.") {
			return global.NULL, err
		} else {
			logger.Println.Error(fmt.Sprintf("获取 %v 存储桶的标签信息时报错，详细信息如下：", bucketName))
			logger.Println.Error(err.Error())
			return global.NULL, err
		}
	} else {
		jsonBytes, err := json.Marshal(response.TagSet)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("解析 %v 存储桶的标签信息时报错，详细信息如下：", bucketName))
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			return string(jsonBytes), err
		}
	}
}
