package qiniuCloud

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/wgpsec/cloudsword/utils/logger"
)

func listKodoBuckets() ([]string, error) {
	var (
		err       error
		buckets   []string
		bucketMgr *storage.BucketManager
	)

	// 获取 Kodo 客户端
	bucketMgr, err = KodoClient()
	if err != nil {
		logger.Println.Error("创建 Kodo 客户端失败，详细信息如下：")
		logger.Println.Error(err.Error())
		fmt.Println()
		return buckets, err
	}

	// 获取所有存储空间（buckets）
	bucketsResp, err := bucketMgr.Buckets(true) // 获取所有存储空间，包括私有和共享
	if err != nil {
		logger.Println.Error("列出存储桶时报错，详细信息如下：")
		logger.Println.Error(err.Error())
		fmt.Println()
	} else {
		buckets = append(buckets, bucketsResp...)
	}

	return buckets, err
}
