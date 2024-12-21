package aliyun

import (
	"fmt"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// ram_create_access_key

func RAMCreateAccessKey() {
	client, err := RAMClient()
	if err == nil {
		userName := global.GetBasicOptionValue(global.UserName)
		createAccessKeyRequest := &ram.CreateAccessKeyRequest{}
		createAccessKeyRequest.UserName = tea.String(userName)
		createAccessKeyResponse, err := client.CreateAccessKey(createAccessKeyRequest)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("为 %v 用户创建访问凭证时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
		} else {
			logger.Println.Info(fmt.Sprintf("为 %v 用户创建访问凭证如下：", userName))
			fmt.Printf("AccessKeyId: %v\nAccessKeySecret: %v\n", *createAccessKeyResponse.Body.AccessKey.AccessKeyId,
				*createAccessKeyResponse.Body.AccessKey.AccessKeySecret)
		}
	}
}
