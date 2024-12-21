package huaweiCloud

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// OBS

func OBSClient() (*obs.ObsClient, error) {
	var (
		err       error
		obsClient *obs.ObsClient
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken == global.NULL {
		obsClient, err = obs.New(AKId, AKSecret, "https://obs."+global.HuaweiCloudDefaultRegion+".myhuaweicloud.com")
	} else {
		obsClient, err = obs.New(AKId, AKSecret, "https://obs."+global.HuaweiCloudDefaultRegion+".myhuaweicloud.com", obs.WithSecurityToken(AKToken))
	}
	if err != nil {
		logger.Println.Error("创建 OBS 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return obsClient, err
}
