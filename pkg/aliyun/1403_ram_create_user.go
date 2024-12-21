package aliyun

import (
	"fmt"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// ram_create_user

func RAMCreateUser() {
	client, err := RAMClient()
	if err == nil {
		userName := global.GetBasicOptionValue(global.UserName)
		createUserRequest := &ram.CreateUserRequest{}
		createUserRequest.UserName = tea.String(userName)
		createUserRequest.Comments = tea.String(global.GetBasicOptionValue(global.Description))
		_, err := client.CreateUser(createUserRequest)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("创建 %v 用户时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
		} else {
			logger.Println.Info(fmt.Sprintf("%v 用户创建成功。", userName))
		}
	}
}
