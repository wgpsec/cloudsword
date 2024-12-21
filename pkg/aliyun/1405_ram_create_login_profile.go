package aliyun

import (
	"fmt"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// ram_create_login_profile

func RAMCreateLoginProfile() {
	client, err := RAMClient()
	if err == nil {
		userName := global.GetBasicOptionValue(global.UserName)
		password := utils.GenerateRandomPasswords()
		createLoginProfileRequest := &ram.CreateLoginProfileRequest{}
		createLoginProfileRequest.UserName = tea.String(userName)
		createLoginProfileRequest.Password = tea.String(password)
		_, err = client.CreateLoginProfile(createLoginProfileRequest)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("为 %v 用户创建 Web 控制台登录配置时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
		} else {
			alias, err := RAMGetAccountAlias()
			if err == nil {
				logger.Println.Info(fmt.Sprintf("为 %v 用户创建 Web 控制台登录配置成功，登录信息如下：", userName))
				fmt.Printf("用户名：%s@%s\n密码：%s\n登录地址：https://signin.aliyun.com\n", userName, alias, password)
			}
		}
	}
}
