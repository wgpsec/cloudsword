package tencentCloud

import (
	"github.com/wgpsec/cloudsword/utils/global"
)

// cam_create_user

func CAMCreateUser() {
	userName := global.GetBasicOptionValue(global.UserName)
	remark := global.GetBasicOptionValue(global.Description)
	_ = addUser(userName, remark)
}
