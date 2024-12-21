package tencentCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// cam_attach_policy_to_user

func CAMAttachPolicyToUser() {
	userName := global.GetBasicOptionValue(global.UserName)
	policyName := global.GetBasicOptionValue(global.PolicyName)
	getUserResponseParams, err := getUser(userName)
	if err == nil {
		strategyInfo, err := listPolicies()
		if err == nil {
			var (
				uin      uint64
				policyId uint64
			)

			for _, policy := range strategyInfo {
				if policyName == *policy.PolicyName {
					policyId = *policy.PolicyId
				}
			}
			if policyId == 0 {
				logger.Println.Error(fmt.Sprintf("预设策略中不存在 %v 策略", policyName))
			} else {
				uin = *getUserResponseParams.Uin
				_ = attachUserPolicy(userName, policyName, uin, policyId)
			}
		}
	}
}
