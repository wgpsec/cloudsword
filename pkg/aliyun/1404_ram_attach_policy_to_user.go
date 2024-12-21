package aliyun

import (
	"fmt"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// ram_attach_policy_to_user

func RAMAttachPolicyToUser() {
	client, err := RAMClient()
	if err == nil {
		userName := global.GetBasicOptionValue(global.UserName)
		policyName := global.GetBasicOptionValue(global.PolicyName)
		attachPolicyToUserRequest := &ram.AttachPolicyToUserRequest{}
		attachPolicyToUserRequest.UserName = tea.String(userName)
		attachPolicyToUserRequest.PolicyName = tea.String(policyName)
		attachPolicyToUserRequest.PolicyType = tea.String("System")
		_, err = client.AttachPolicyToUser(attachPolicyToUserRequest)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("为 %v 用户添加 %v 策略时报错，详细信息如下：", userName, policyName))
			logger.Println.Error(err.Error())
		} else {
			logger.Println.Info(fmt.Sprintf("为 %v 用户添加 %v 策略成功。", userName, policyName))
		}
	}
}
