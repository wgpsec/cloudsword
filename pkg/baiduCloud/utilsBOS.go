package baiduCloud

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/wgpsec/cloudsword/utils/logger"
)

func listBOSBuckets() ([]api.BucketSummaryType, error) {
	var (
		err      error
		client   *bos.Client
		response *api.ListBucketsResult
		buckets  []api.BucketSummaryType
	)
	client, err = BOSClient()
	if err == nil {
		response, err = client.ListBuckets()
		if err != nil {
			logger.Println.Error("列出存储桶时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			fmt.Println()
		} else {
			buckets = append(buckets, response.Buckets...)
		}
	}
	return buckets, err
}
