package tencentCloud

import (
	"encoding/json"
	"fmt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// cam_list_roles

func CAMListRoles() {
	roles, err := getCAMRoles()
	if err == nil {
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
					utils.PrintfNotNilString("角色类型：", role.RoleType)
					utils.PrintfNotNilString("角色资源描述符 ARN：", role.RoleArn)
					utils.PrintfNotNilUInt64("角色最大会话时间：", role.SessionDuration)
					fmt.Printf("角色描述：%v\n", *role.Description)
					fmt.Printf("角色策略文档：%v\n", *role.PolicyDocument)
					jsonBytes, err := json.Marshal(role.Tags)
					if err == nil {
						fmt.Printf("实例标签：%s\n", jsonBytes)
					}
					if role.AddTime == nil {
						fmt.Printf("创建时间：%v\n", *role.AddTime)
					}
					if role.UpdateTime == nil {
						fmt.Printf("更新时间：%v\n", *role.UpdateTime)
					}
				}
			}
			fmt.Println()
		}
	}
}
