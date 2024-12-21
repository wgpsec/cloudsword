package aliyun

import (
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// 获取账号 ID

func RAMGetAccountAlias() (string, error) {
	client, err := RAMClient()
	if err == nil {
		GetAccountAliasResponse, err := client.GetAccountAlias()
		if err != nil {
			logger.Println.Error("查看云账号别名时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			return *GetAccountAliasResponse.Body.AccountAlias, err
		}
	} else {
		return global.NULL, err
	}
}
