package aliyun

import (
	"fmt"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// ram_list_users

func RAMListUsers() {
	var users []*ram.ListUsersResponseBodyUsersUser
	listUsersRequest := &ram.ListUsersRequest{}
	listUsersRequest.MaxItems = tea.Int32(1000)
	runtime := &util.RuntimeOptions{}

	for {
		client, err := RAMClient()
		if err == nil {
			response, err := client.ListUsersWithOptions(listUsersRequest, runtime)
			if err != nil {
				logger.Println.Error("列出 RAM 用户时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				break
			} else {
				users = append(users, response.Body.Users.User...)
				if *response.Body.IsTruncated {
					listUsersRequest.Marker = response.Body.Marker
				} else {
					break
				}
			}
		}
	}
	if len(users) == 0 {
		logger.Println.Info("未找到用户。\n")
	} else {
		logger.Println.Info("找到以下用户：")
		for _, user := range users {
			if global.GetBasicOptionValue(global.Detail) == global.False {
				fmt.Println(*user.UserName)
			} else {
				fmt.Printf("\n用户登录名：%v\n", *user.UserName)
				fmt.Printf("用户显示名称：%v\n", *user.DisplayName)
				fmt.Printf("用户 ID：%v\n", *user.UserId)
				fmt.Printf("邮箱：%v\n", utils.ConvertedNullPointer(user.Email))
				fmt.Printf("手机号：%v\n", utils.ConvertedNullPointer(user.MobilePhone))
				fmt.Printf("备注：%v\n", utils.ConvertedNullPointer(user.Comments))

				createDate, err := utils.GetUTC8TimeType2(*user.CreateDate)
				if err == nil {
					fmt.Printf("创建时间：%v\n", createDate)
				}
			}
		}
		fmt.Println()
	}
}
