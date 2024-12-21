package huaweiCloud

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/wgpsec/cloudsword/utils/logger"
)

func listOBSBuckets() ([]obs.Bucket, error) {
	var (
		err     error
		client  *obs.ObsClient
		resp    *obs.ListBucketsOutput
		buckets []obs.Bucket
	)
	client, err = OBSClient()
	if err == nil {
		obsListBucketsInput := &obs.ListBucketsInput{}
		obsListBucketsInput.QueryLocation = true
		resp, err = client.ListBuckets(obsListBucketsInput)
		if err != nil {
			logger.Println.Error("列出存储桶时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			fmt.Println()
		} else {
			buckets = append(buckets, resp.Buckets...)
		}
	}
	return buckets, err
}
