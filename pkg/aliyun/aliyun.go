package aliyun

import (
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
)

func Modules() []global.Module {
	var modules []global.Module
	modules = append(modules, modules1()...)
	modules = append(modules, modules2()...)
	modules = append(modules, modules3()...)
	modules = append(modules, modules4()...)
	modules = append(modules, modules5()...)
	return modules
}

// 综合性模块

func modules1() []global.Module {
	return []global.Module{
		utils.GetModule(
			1101,
			4,
			global.Aliyun,
			"list_cloud_assets",
			global.TeamsSix,
			"列出 OSS、ECS、RAM、Domain 服务资产",
			`
简洁输出所需权限：
1. oss:ListBuckets
2. ecs:DescribeRegions
3. ecs:DescribeInstances
4. ram:ListUsers
5. ram:ListRoles
6. domain:QueryDomainList

详细输出所需权限：
1. oss:ListBuckets
2. oss:GetBucketStat
3. oss:GetBucketTags
4. ecs:DescribeRegions
5. ecs:DescribeInstances
6. ecs:DescribeCloudAssistantStatus
7. ecs:DescribeUserData
8. ram:ListUsers
9. ram:ListRoles
10. domain:QueryDomainList
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail, global.BasicOptionRegion),
		),
	}
}

// OSS 模块

func modules2() []global.Module {
	return []global.Module{
		utils.GetModule(
			1201,
			2,
			global.Aliyun,
			"oss_list_buckets",
			global.TeamsSix,
			"列出阿里云 OSS 对象存储桶",
			`
简洁输出模式所需权限：
1. oss:ListBuckets

详细输出模式所需权限：
1. oss:ListBuckets
2. oss:GetBucketStat
3. oss:GetBucketTags
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),

		utils.GetModule(
			1202,
			4,
			global.Aliyun,
			"oss_search_objects",
			global.TeamsSix,
			"搜索阿里云 OSS 对象",
			`
所需权限：
1. oss:GetBucketLocation
2. oss:GetMetaQueryStatus
3. oss:OpenMetaQuery
4. oss:DoMetaQuery
`,
			append(
				global.BasicOptionsDefault,
				global.BasicOptionDetail,
				global.BasicOptionQueryValue,
				global.BasicOptionMaxResults,
				global.BasicOptionBucketName,
			),
		),

		utils.GetModule(
			1203,
			3,
			global.Aliyun,
			"oss_bucket_only_upload_images",
			global.TeamsSix,
			"使用云函数限制存储桶只允许上传图片",
			`
所需权限：
1. ram:CreateRole
2. ram:CreatePolicy
3. ram:AttachPolicyToRole
4. fc:CreateFunction
5. fc:CreateTrigger

使用说明：
1. 使用该功能前请仔细阅读《使用云函数限制存储桶上传类型》这篇文章，文章地址：https://wiki.teamssix.com/cloudservice/s3/limiting-bucket-upload-types-using-cloud-functions.html
2. 这个功能的运行逻辑是文件会先被上传到存储桶中，然后云函数判断刚刚上传的文件是否是预期类型，如果不是就会将这个文件删除。
3. 这个功能设置后，所指定的存储桶只能上传 jpg、jpeg、png、gif 的文件，如果上传文件的文件头类型、文件后缀类型、Content Type 三个有任何一个不符合预期的都会被直接删除。
4. 这个功能没办法解决在上传同名称文件时被覆盖的问题。
`,
			append(
				global.BasicOptionsDefault,
				global.BasicOptionBucketName,
			),
		),
	}
}

// ECS 模块

func modules3() []global.Module {
	return []global.Module{
		utils.GetModule(
			1301,
			2,
			global.Aliyun,
			"ecs_list_instances",
			global.TeamsSix,
			"列出阿里云 ECS 弹性计算实例",
			`
简洁输出模式所需权限：
1. ecs:DescribeRegions
2. ecs:DescribeInstances

详细输出模式所需权限：
1. ecs:DescribeRegions
2. ecs:DescribeInstances
3. ecs:DescribeCloudAssistantStatus
4. ecs:DescribeUserData
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail, global.BasicOptionRegion),
		),
	}
}

// RAM 模块

func modules4() []global.Module {
	return []global.Module{

		utils.GetModule(
			1401,
			2,
			global.Aliyun,
			"ram_list_users",
			global.TeamsSix,
			"列出阿里云 RAM 用户",
			`
所需权限：
1. ram:ListUsers
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
		utils.GetModule(
			1402,
			1,
			global.Aliyun,
			"ram_list_roles",
			global.TeamsSix,
			"列出阿里云 RAM 角色",
			`
所需权限：
1. ram:ListRoles
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
		utils.GetModule(
			1403,
			1,
			global.Aliyun,
			"ram_create_user",
			global.TeamsSix,
			"创建阿里云 RAM 用户",
			`
所需权限：
1. ram:CreateUser
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName, global.BasicOptionDescription),
		),
		utils.GetModule(
			1404,
			1,
			global.Aliyun,
			"ram_attach_policy_to_user",
			global.TeamsSix,
			"为阿里云 RAM 用户添加策略",
			`
所需权限：
1. ram:AttachPolicyToUser

使用频率较高的策略名称（仅供参考）：
1. AdministratorAccess
2. ReadOnlyAccess
3. AliyunRAMFullAccess
4. AliyunOSSFullAccess
5. AliyunECSFullAccess

注意事项：
1. 为用户添加 AdministratorAccess 或 AliyunRAMFullAccess 策略是一件高危操作，请根据实际情况酌量进行操作，此操作存在引发告警的风险。
2. 此功能仅支持添加系统策略，自定义策略建议使用控制台或 aliyun cli 等工具进行添加，完整的权限策略清单可参考：https://ram.console.aliyun.com/policies
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName, global.BasicOptionPolicyName),
		),
		utils.GetModule(
			1405,
			3,
			global.Aliyun,
			"ram_create_login_profile",
			global.TeamsSix,
			"创建阿里云 RAM 用户 Web 登录配置",
			`
所需权限：
1. ram:CreateLoginProfile
2. ram:GetAccountAlias
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName),
		),
		utils.GetModule(
			1406,
			1,
			global.Aliyun,
			"ram_create_access_key",
			global.TeamsSix,
			"创建阿里云 RAM 用户访问凭证",
			`
所需权限：
1. ram:CreateAccessKey
`,
			append(global.BasicOptionsDefault, global.BasicOptionUserName),
		),
	}
}

// 综合性模块

func modules5() []global.Module {
	return []global.Module{
		// 其他模块
		utils.GetModule(
			1501,
			1,
			global.Aliyun,
			"domain_list_domains",
			global.TeamsSix,
			"列出阿里云 Domains 域名资产",
			`
所需权限：
1. domain:QueryDomainList
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
	}
}
