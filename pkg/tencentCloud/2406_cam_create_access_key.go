package tencentCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// cam_create_access_key

func CAMCreateAccessKey() {
	userName := global.GetBasicOptionValue(global.UserName)
	getUserResponseParams, err := getUser(userName)
	if err == nil {
		accessKeyDetail, err := createAccessKey(userName, *getUserResponseParams.Uin)
		if err == nil {
			logger.Println.Info(fmt.Sprintf("为 %v 用户创建访问凭证如下：", userName))
			fmt.Printf("AccessKeyId: %v\nAccessKeySecret: %v\n", *accessKeyDetail.AccessKeyId,
				*accessKeyDetail.SecretAccessKey)
		}
	}
}
