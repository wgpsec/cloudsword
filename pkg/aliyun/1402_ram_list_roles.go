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

// ram_list_roles

func RAMListRoles() {
	var roles []*ram.ListRolesResponseBodyRolesRole
	listRolesRequest := &ram.ListRolesRequest{}
	listRolesRequest.MaxItems = tea.Int32(1000)
	runtime := &util.RuntimeOptions{}

	for {
		client, err := RAMClient()
		if err == nil {
			response, err := client.ListRolesWithOptions(listRolesRequest, runtime)
			if err != nil {
				logger.Println.Error("列出 RAM 角色时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				break
			} else {
				roles = append(roles, response.Body.Roles.Role...)
				if *response.Body.IsTruncated {
					listRolesRequest.Marker = response.Body.Marker
				} else {
					break
				}
			}
		}
	}
	if len(roles) == 0 {
		logger.Println.Info("未找到角色。\n")
	} else {
		logger.Println.Info("找到以下角色：")
		for _, role := range roles {
			if global.GetBasicOptionValue(global.Detail) == global.False {
				fmt.Println(*role.RoleName)
			} else {
				fmt.Printf("\n角色名称：%v\n", *role.RoleName)
				fmt.Printf("角色 ID：%v\n", *role.RoleId)
				fmt.Printf("角色资源描述符 ARN：%v\n", *role.Arn)
				fmt.Printf("角色最大会话时间：%v\n", *role.MaxSessionDuration)
				fmt.Printf("角色描述：%v\n", *role.Description)

				createDate, err := utils.GetUTC8TimeType2(*role.CreateDate)
				if err == nil {
					fmt.Printf("创建时间：%v\n", createDate)
				}

				updateDate, err := utils.GetUTC8TimeType2(*role.UpdateDate)
				if err == nil {
					fmt.Printf("更新时间：%v\n", updateDate)
				}
			}
		}
		fmt.Println()
	}
}
