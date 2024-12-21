package tencentCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// cam_create_login_profile

func CAMCreateLoginProfile() {
	userName := global.GetBasicOptionValue(global.UserName)
	password := utils.GenerateRandomPasswords()
	err := updateUserWithConsoleLogin(userName, password)
	if err == nil {
		appIdResponse, err := getUserAppId()
		if err == nil {
			logger.Println.Info(fmt.Sprintf("为 %v 用户创建 Web 控制台登录配置成功，登录信息如下：", userName))
			fmt.Printf("主账号 ID：%s\n子用户名：%s\n密码：%s\n登录地址：https://cloud.tencent.com/login/subAccount\n", *appIdResponse.OwnerUin, userName, password)
		}
	}
}
