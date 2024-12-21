package baiduCloud

import (
	"github.com/baidubce/bce-sdk-go/auth"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// BOS

func BOSClient() (*bos.Client, error) {
	var (
		bosClient     *bos.Client
		err           error
		stsCredential *auth.BceCredentials
		endpoint      = "https://bj.bcebos.com"
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken != global.NULL {
		bosClient, err = bos.NewClient(AKId, AKSecret, "")
		stsCredential, err = auth.NewSessionBceCredentials(AKId, AKSecret, AKToken)
		bosClient.Config.Credentials = stsCredential
	} else {
		clientConfig := bos.BosClientConfiguration{
			Ak:               AKId,
			Sk:               AKSecret,
			Endpoint:         endpoint,
			RedirectDisabled: false,
		}
		bosClient, err = bos.NewClientWithConfig(&clientConfig)
	}
	if err != nil {
		logger.Println.Error("创建 BOS 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return bosClient, err
}
