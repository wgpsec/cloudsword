package tencentCloud

import (
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
)

func Modules() []global.Module {
	var modules []global.Module
	modules = append(modules, module1()...)
	modules = append(modules, module2()...)
	modules = append(modules, module3()...)
	modules = append(modules, module4()...)
	return modules
}

// 综合性模块

func module1() []global.Module {
	return []global.Module{
		utils.GetModule(2101,
			4,
			global.TencentCloud,
			"list_cloud_assets",
			global.TeamsSix,
			"列出 COS、EVM、LH、RAM 服务资产",
			`
简洁输出所需权限：
1. cos:GetService
2. cvm:DescribeRegions
3. cvm:DescribeInstances
4. lh:DescribeRegions
5. lh:DescribeInstances
6. cam:ListUsers
7. cam:DescribeRoleList

详细输出所需权限：
1. cos:GetService
2. cos:GetTagging
3. cvm:DescribeRegions
4. cvm:DescribeInstances
5. lh:DescribeRegions
6. lh:DescribeInstances
7. cam:ListUsers
8. cam:DescribeRoleList
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail, global.BasicOptionRegion),
		), utils.GetModule(2102,
			5,
			global.TencentCloud,
			"create_honey_token",
			global.TeamsSix,
			"创建云访问凭证蜜标",
			`
所需权限：
1. cam:AddUser
2. cam:GetUser
3. cam:CreateAccessKey
4. cloudaudit:CreateAuditTrack
5. cls:CreateTopic
6. cls:CreateNoticeContent
7. cls:CreateAlarmNotice
8. cls:CreateAlarm

Webhook 支持的类型有企业微信、钉钉、飞书、自定义接口，各平台 Webhook 生成方式文档地址如下：
1. 企业微信：https://developer.work.weixin.qq.com/document/path/91770
2. 钉钉：https://open.dingtalk.com/document/orgapp/obtain-the-webhook-address-of-a-custom-robot
3. 飞书：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot

使用说明：
1. 使用该功能前请仔细阅读《使用云访问凭证蜜标及时发现入侵行为》这篇文章，文章地址：https://wiki.teamssix.com/cloudservice/more/use-honeytokens-to-discover-attackers.html
2. 这个功能会创建一个访问凭证蜜标，您可以将创建的访问凭证蜜标放到服务器、工作电脑、环境变量、程序代码、配置文件等易于发现的位置，然后当这个访问凭证被调用时，说明访问凭证已经被泄露了，此时就会触发告警，这样您可以及时知道自己被入侵了。
3. 这个功能所返回的访问凭证是没有添加任何权限的，但是为了避免一些意外情况，建议在非生产环境下使用此功能。
4. 这个功能初始化需要 10 分钟左右的时间，当初始化完成后，从蜜标凭证被调用到接收到告警大概会经历 5 分钟左右的时间。
`,
			append(global.BasicOptionsDefault, global.BasicOptionWebhook),
		),
	}
}

// COS 模块

func module2() []global.Module {
	return []global.Module{
		utils.GetModule(
			2201,
			2,
			global.TencentCloud,
			"cos_list_buckets",
			global.TeamsSix,
			"列出腾讯云 COS 对象存储桶",
			`
简洁输出模式所需权限：
1. cos:GetService

详细输出模式所需权限：
1. cos:GetService
2. cos:GetTagging
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
	}
}

// CVM 模块

func module3() []global.Module {
	return []global.Module{
		utils.GetModule(
			2301,
			2,
			global.TencentCloud,
			"cvm_list_instances",
			global.TeamsSix,
			"列出腾讯云 CVM 弹性计算实例",
			`
所需权限：
1. cvm:DescribeRegions
2. cvm:DescribeInstances
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail, global.BasicOptionRegion),
		),
		utils.GetModule(
			2302,
			1,
			global.TencentCloud,
			"lh_list_instances",
			global.TeamsSix,
			"列出腾讯云 LH 轻量应用服务器",
			`
所需权限：
1. lh:DescribeRegions
2. lh:DescribeInstances
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail, global.BasicOptionRegion),
		),
	}
}

// CAM 模块

func module4() []global.Module {
	return []global.Module{
		utils.GetModule(
			2401,
			2,
			global.TencentCloud,
			"cam_list_users",
			global.TeamsSix,
			"列出腾讯云 CAM 用户",
			`
所需权限：
1. cam:ListUsers
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
		utils.GetModule(
			2402,
			1,
			global.TencentCloud,
			"cam_list_roles",
			global.TeamsSix,
			"列出腾讯云 CAM 角色",
			`
所需权限：
1. cam:DescribeRoleList
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
		utils.GetModule(
			2403,
			1,
			global.TencentCloud,
			"cam_create_user",
			global.TeamsSix,
			"创建腾讯云 CAM 角色",
			`
所需权限：
1. cam:AddUser
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName, global.BasicOptionDescription),
		),
		utils.GetModule(
			2404,
			1,
			global.TencentCloud,
			"cam_attach_policy_to_user",
			global.TeamsSix,
			"为腾讯云 CAM 用户添加策略",
			`
所需权限：
1. cam:GetUser
2. cam:ListPolicies
3. cam:AttachUserPolicy

使用频率较高的策略名称（仅供参考）：
1. AdministratorAccess
2. ReadOnlyAccess
3. QcloudCamFullAccess
4. QcloudCOSFullAccess
5. QcloudCVMFullAccess

注意事项：
1. 为用户添加 AdministratorAccess 或 QcloudCamFullAccess 策略是一件高危操作，请根据实际情况酌量进行操作。
2. 此功能仅支持添加预设策略，自定义策略建议使用控制台或 TCCLI 等工具进行添加，完整的权限策略清单可参考：https://console.cloud.tencent.com/cam/policy
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName, global.BasicOptionPolicyName),
		),
		utils.GetModule(
			2405,
			3,
			global.TencentCloud,
			"cam_create_login_profile",
			global.TeamsSix,
			"创建腾讯云 CAM 用户 Web 登录配置",
			`
所需权限：
1. cam:UpdateUser
2. cam:GetUserAppId
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName),
		),
		utils.GetModule(
			2406,
			1,
			global.TencentCloud,
			"cam_create_access_key",
			global.TeamsSix,
			"创建腾讯云 CAM 用户访问凭证",
			`
所需权限：
1. cam:GetUser
2. cam:CreateAccessKey
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName),
		),
	}
}
