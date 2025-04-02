package qiniuCloud

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

//Kodo

func KodoClient() (*storage.BucketManager, error) {
	var (
		err       error
		mac       *qbox.Mac
		cfg       storage.Config
		bucketMgr *storage.BucketManager
	)

	// 获取配置信息
	accessKey := global.GetBasicOptionValue(global.AKId)
	secretKey := global.GetBasicOptionValue(global.AKSecret)

	// 创建认证信息
	mac = qbox.NewMac(accessKey, secretKey)

	// 配置存储参数
	cfg = storage.Config{
		UseHTTPS: true, // 根据需要选择是否使用 HTTPS
	}

	// 创建 BucketManager 实例
	bucketMgr = storage.NewBucketManager(mac, &cfg)

	if bucketMgr == nil {
		err = fmt.Errorf("创建 Kodo 客户端失败")
		logger.Println.Error("创建 Kodo 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}

	return bucketMgr, err
}
