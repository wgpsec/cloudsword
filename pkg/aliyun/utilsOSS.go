package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

// 列存储桶

func listBuckets() []oss.BucketProperties {
	var buckets []oss.BucketProperties
	listBucketsRequest := &oss.ListBucketsRequest{}
	for {
		response, err := OSSClient(global.AliyunDefaultRegion).ListBuckets(context.TODO(), listBucketsRequest)
		if err != nil {
			logger.Println.Error("列出存储桶时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			fmt.Println()
			break
		} else {
			if response.IsTruncated {
				listBucketsRequest.Marker = response.NextMarker
			} else {
				buckets = append(buckets, response.Buckets...)
				break
			}
		}
	}
	return buckets
}

// 查询 Bucket 容量信息

func getBucketStat(bucketName, region string) (*oss.GetBucketStatResult, error) {
	getBucketStatRequest := &oss.GetBucketStatRequest{}
	getBucketStatRequest.Bucket = tea.String(bucketName)
	response, err := OSSClient(region).GetBucketStat(context.TODO(), getBucketStatRequest)
	if err != nil {
		logger.Println.Error(fmt.Sprintf("获取 %v 存储桶的容量信息时报错，详细信息如下：", bucketName))
		logger.Println.Error(err.Error())
	}
	return response, err
}

// 查询 Bucket 标签信息

func getBucketTags(BucketName, region string) (string, error) {
	getBucketTagsRequest := &oss.GetBucketTagsRequest{}
	getBucketTagsRequest.Bucket = tea.String(BucketName)
	response, err := OSSClient(region).GetBucketTags(context.TODO(), getBucketTagsRequest)
	if err != nil {
		logger.Println.Error(fmt.Sprintf("获取 %v 存储桶的标签信息时报错，详细信息如下：", BucketName))
		logger.Println.Error(err.Error())
		return global.NULL, err
	} else {
		jsonBytes, err := json.Marshal(response.Tagging.TagSet.Tags)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("解析 %v 存储桶的标签信息时报错，详细信息如下：", BucketName))
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			return string(jsonBytes), err
		}
	}
}

// 查询 Bucket 数据索引状态

func getMetaQueryStatus(bucketName, region string) (string, error) {
	getMetaQueryStatusRequest := &oss.GetMetaQueryStatusRequest{}
	getMetaQueryStatusRequest.Bucket = tea.String(bucketName)
	response, err := OSSClient(region).GetMetaQueryStatus(context.TODO(), getMetaQueryStatusRequest)
	if err == nil {
		return *response.MetaQueryStatus.State, err
	} else if strings.Contains(err.Error(), "MetaQueryNotExist") {
		return "MetaQueryNotExist", nil
	} else {
		logger.Println.Error(fmt.Sprintf("获取 %v 存储桶数据索引时报错，详细信息如下：", bucketName))
		logger.Println.Error(err.Error())
		return global.NULL, err
	}
}

// 开启数据索引

func openMetaQuery(bucketName, region string) error {
	openMetaQueryRequest := &oss.OpenMetaQueryRequest{}
	openMetaQueryRequest.Bucket = tea.String(bucketName)
	_, err := OSSClient(region).OpenMetaQuery(context.TODO(), openMetaQueryRequest)
	if err == nil {
		logger.Println.Info("数据索引已开启。")
		return err
	} else if strings.Contains(err.Error(), "MetaQueryAlreadyExist") {
		return nil
	} else {
		logger.Println.Error(fmt.Sprintf("开启 %v 存储桶的数据索引时报错，详细信息如下：", bucketName))
		logger.Println.Error(err.Error())
		return err
	}
}

// 关闭数据索引

func closeMetaQuery(bucketName, region string) {
	closeMetaQueryRequest := &oss.CloseMetaQueryRequest{}
	closeMetaQueryRequest.Bucket = tea.String(bucketName)
	_, err := OSSClient(region).CloseMetaQuery(context.TODO(), closeMetaQueryRequest)
	if err == nil {
		logger.Println.Info("数据索引已关闭。")
	} else {
		logger.Println.Error(fmt.Sprintf("关闭 %v 存储桶的数据索引时报错，详细信息如下：", bucketName))
		logger.Println.Error(err.Error())
	}
}

// 查询存储桶区域

func getBucketRegion(bucketName string) (string, error) {
	getBucketLocationRequest := &oss.GetBucketLocationRequest{}
	getBucketLocationRequest.Bucket = tea.String(bucketName)
	response, err := OSSClient(global.AliyunDefaultRegion).GetBucketLocation(context.TODO(), getBucketLocationRequest)
	if err == nil {
		return strings.Replace(*response.LocationConstraint, "oss-", "", -1), err
	} else if strings.Contains(err.Error(), "OperationInput.Bucket") {
		logger.Println.Error(fmt.Sprintf("%v 存储桶不存在。", bucketName))
		return global.NULL, err
	} else {
		logger.Println.Error(fmt.Sprintf("获取 %v 存储桶的区域时报错，详细信息如下：", bucketName))
		logger.Println.Error(err.Error())
		return global.NULL, err
	}
}
