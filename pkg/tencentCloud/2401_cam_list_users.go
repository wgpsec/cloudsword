package tencentCloud

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// cam_list_users

func CAMListUsers() {
	users, err := getCAMUsers()
	if err == nil {
		if len(users) == 0 {
			logger.Println.Info("未找到用户。\n")
		} else {
			logger.Println.Info("找到以下用户：")
			for _, user := range users {
				if global.GetBasicOptionValue(global.Detail) == global.False {
					fmt.Println(*user.Name)
				} else {
					fmt.Printf("\n用户登录名：%v\n", *user.Name)
					utils.PrintfNotNilString("用户昵称：", user.NickName)
					fmt.Printf("用户 ID：%v\n", *user.Uin)
					fmt.Printf("用户 UID：%v\n", *user.Uid)
					fmt.Printf("邮箱：%v\n", *user.Email)
					fmt.Printf("手机号：%v\n", *user.PhoneNum)
					fmt.Printf("备注：%v\n", *user.Remark)
					if user.CreateTime != nil {
						fmt.Printf("创建时间：%v\n", *user.CreateTime)
					}
				}
			}
			fmt.Println()
		}
	}
}
