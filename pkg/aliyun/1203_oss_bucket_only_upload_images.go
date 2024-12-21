package aliyun

import (
	"fmt"
	fc "github.com/alibabacloud-go/fc-20230330/v4/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"regexp"
	"strings"
)

func OSSLimitBucketOnlyUploadImages() {
	// 1. 创建 RAM 角色
	ramClient, err := RAMClient()
	if err == nil {
		roleName := utils.GenerateRandomName("Role")
		logger.Println.Info(fmt.Sprintf("正在创建 %v 角色。", roleName))
		createRoleRequest := &ram.CreateRoleRequest{}
		createRoleRequest.RoleName = tea.String(roleName)
		createRoleRequest.Description = tea.String("此角色由云鉴 1203_oss_bucket_only_upload_images 模块自动创建，此角色用于云函数读取和删除存储桶对象以及执行云函数使用。")
		createRoleRequest.AssumeRolePolicyDocument = tea.String(`
{
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "fc.aliyuncs.com"
                ]
            }
        },
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "Service": [
                    "oss.aliyuncs.com"
                ]
            }
        }
    ],
    "Version": "1"
}
`)
		createRoleResponse, err := ramClient.CreateRole(createRoleRequest)
		if err != nil {
			logger.Println.Error("创建 RAM 角色时报错，详细信息如下：")
			logger.Println.Error(err.Error())
		} else {
			// 2. 创建权限策略
			roleArn := *createRoleResponse.Body.Role.Arn
			accountId := regexp.MustCompile(`::(\d+):`).FindStringSubmatch(roleArn)[1]
			bucketName := global.GetBasicOptionValue(global.BucketName)
			region, err := getBucketRegion(bucketName)
			if err == nil {
				policyName := utils.GenerateRandomName("Policy")
				logger.Println.Info(fmt.Sprintf("正在创建 %v 策略。", policyName))
				createPolicyRequest := &ram.CreatePolicyRequest{}
				createPolicyRequest.PolicyName = tea.String(policyName)
				createPolicyRequest.Description = tea.String(
					"此策略由云鉴 1203_oss_bucket_only_upload_images 模块自动创建，此策略所关联的角色用于云函数读取和删除存储桶对象以及执行云函数使用。")
				createPolicyRequest.PolicyDocument = tea.String(fmt.Sprintf(`
{
    "Version": "1",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "oss:GetObject",
                "oss:DeleteObject"
            ],
            "Resource": "acs:oss:oss-%s:%s:%s/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "fc:InvokeFunction"
            ],
            "Resource": "*"
        }
    ]
}
`, region, accountId, bucketName))
				_, err := ramClient.CreatePolicy(createPolicyRequest)
				if err != nil {
					logger.Println.Error("创建 RAM 策略时报错，详细信息如下：")
					logger.Println.Error(err.Error())
					logger.Println.Info(fmt.Sprintf(`
已经创建的资源如下：
角色：%v
`, roleName))
				} else {
					// 3. 为 RAM 角色授予权限
					logger.Println.Info(fmt.Sprintf("正在将 %v 策略绑定到 %v 角色。", policyName, roleName))
					attachPolicyToRoleRequest := &ram.AttachPolicyToRoleRequest{}
					attachPolicyToRoleRequest.PolicyName = tea.String(policyName)
					attachPolicyToRoleRequest.RoleName = tea.String(roleName)
					attachPolicyToRoleRequest.PolicyType = tea.String("Custom")

					_, err = ramClient.AttachPolicyToRole(attachPolicyToRoleRequest)
					if err != nil {
						logger.Println.Error(fmt.Sprintf("为 %v 角色绑定 %v 策略时报错，详细信息如下：", roleName, policyName))
						logger.Println.Error(err.Error())
						logger.Println.Info(fmt.Sprintf(`
已经创建的资源如下：
角色：%v
策略：%v
`, roleName, policyName))
					} else {
						// 4. 创建云函数
						fcClient, err := FCClient(accountId, region)
						if err == nil {
							functionName := utils.GenerateRandomName("Function")
							logger.Println.Info(fmt.Sprintf("正在创建 %v 函数。", functionName))
							createFunctionRequest := &fc.CreateFunctionRequest{}
							createFunctionRequest.Body = &fc.CreateFunctionInput{}
							createFunctionRequest.Body.FunctionName = tea.String(functionName)
							createFunctionRequest.Body.Runtime = tea.String("python3.10")
							createFunctionRequest.Body.Handler = tea.String("index.handler")
							createFunctionRequest.Body.Cpu = tea.Float32(0.35)
							createFunctionRequest.Body.MemorySize = tea.Int32(512)
							createFunctionRequest.Body.DiskSize = tea.Int32(512)
							createFunctionRequest.Body.Role = tea.String(roleArn)
							createFunctionRequest.Body.Timeout = tea.Int32(60)
							createFunctionRequest.Body.Description = tea.String(
								"此函数由云鉴 1203_oss_bucket_only_upload_images 模块自动创建，此函数用于识别上传到存储桶里的对象是否是图片，如果不是则删除。")
							createFunctionRequest.Body.Code = &fc.InputCodeLocation{}
							createFunctionRequest.Body.Code.ZipFile = tea.String("UEsDBBQACAAIANS6hVkAAAAAAAAAAP0JAAAIACAAaW5kZXgucHlVVA0AB0HFUWcJz1NnB89TZ3V4CwABBPUBAAAEFAAAAI1VTU/bQBC951e46mETERygNyQOLfde2lsURcaeJAbHtrwbPoSQQgUFtZRyghaqFgqlSBWJ2iI+RFT+S4uDOfEXut61w8axgSiSdz0zb968nVnrVdtyiGThlM5X49gyw7WF8VC41qvliuakUikNSlJFMTUDnDRMgkmykmqZBKZJZjgl0d9jyVs9cd+vX54etTeOvcbF9UbDXTr0mvPu2bF3seXtrHiNPbpl3qoDGpZGQgzZ31NQXTEwsys1UqFmn4r8guCndJtmMbKiqoBxcQJmirqWlXpeYqCvSGigu5qjk5kisSbAzKQiVN9ctrbb60uX58dXmwtXW6fu4u/Li532fJP5wSQpGphQHr48smEpGubFZ0I7tQVeecQsGBXyAwVmHqupE0CKplIF7pZHtB5qR9zir3wj4u7W2DioxK8h4s0N/oraAmcwNdvSTT+/79aPpD4e40BZt0xUoHvUTx3AMRVDVgx9pmaqWFatKhLIhRI/Y7u0r3q2g50VK7hbOndt9apVv/p57n5+y/xqtq9WsaQbUMS1UkmfLpIZG1g+2VZIRca2oftHn76tO5MfLMgO2IaiQhrJKCshdE/e14tu40zIWwGa1QFcM/zaOH+ZveRpxGw9RFkzmiRkKmAxCHCwXKYqoVHu1/+S+t1DkC/cvaMEbYKUt1xpghiqVBVFS/cy5rRCwnxY5amKQtLPLROyUmUkJll++MlQIaBtO/Sg0yWU7ytInCsvwDtpun8WblqbtD2Gc7lZoRPmcrO3zOZQJgGnUzMFmU3gfEe00FBRAKGfkgEk8ZCiCOJBz92eoLs47zVOuQCds+Mc3OUNb+eAX02GYU2BxoL9KyyPxu2y36vjNrCnbbJHWS8FwxpGiGl5pF5VypALA/kuCOcbDhLwa+/W21++ibTaH5ru2r67tux9X2BOo1aVzpXCOmpAkIYrc7jrfjrgKJEqRSDpb32f/gNp9VJix5kW/UCY3YLwj4GYuL+nI0IG767rm97F0r/6qyAXGBh6AXpbigNEo2PYilfP3WzjGQtdmEi6l3g8+R6sOKCYGrpupUgRXR3VTSC2oHAqJH8qpDtLii8rtrQ41CTIaJMOdoZvefv64x4lEzuFYSt2QqNt0jX7911aCRzj2q8kqvdA3I6iN60VPnO8OK8ZDG+3LMHNr4EBBJK+U12MHj2ckXvyy/v6g+fnaVP/AVBLBwhBqwo4hAMAAP0JAABQSwECFAMUAAgACADUuoVZQasKOIQDAAD9CQAACAAgAAAAAAAAAAAApIEAAAAAaW5kZXgucHlVVA0AB0HFUWcJz1NnB89TZ3V4CwABBPUBAAAEFAAAAFBLBQYAAAAAAQABAFYAAADaAwAAAAA=")
							createFunctionResponse, err := fcClient.CreateFunction(createFunctionRequest)
							if err != nil {
								logger.Println.Error("创建函数时报错，详细信息如下：")
								logger.Println.Error(err.Error())
								logger.Println.Info(fmt.Sprintf(`
已经创建的资源如下：
角色：%v
策略：%v
`, roleName, policyName))
							} else {
								// 5. 创建触发器
								triggerName := utils.GenerateRandomName("Trigger")
								logger.Println.Info(fmt.Sprintf("正在创建 %v 触发器。", triggerName))
								createTriggerRequest := &fc.CreateTriggerRequest{}
								createTriggerRequest.Body = &fc.CreateTriggerInput{}
								createTriggerRequest.Body.Description = tea.String(
									"此触发器由云鉴 1203_oss_bucket_only_upload_images 模块自动创建，此触发器用于识别当前是否存在存储桶文件上传动作。")
								createTriggerRequest.Body.TriggerType = tea.String("oss")
								createTriggerRequest.Body.TriggerName = tea.String(triggerName)
								createTriggerRequest.Body.TriggerConfig = tea.String(`{"events":["oss:ObjectCreated:PutObject",
"oss:ObjectCreated:PostObject","oss:ObjectCreated:CompleteMultipartUpload","oss:ObjectCreated:PutSymlink"],
"filter":{"key":{"prefix":"","suffix":""}}}`)
								createTriggerRequest.Body.InvocationRole = tea.String(roleArn)
								createTriggerRequest.Body.SourceArn = tea.String(fmt.Sprintf("acs:oss:%s:%s:%s", region, accountId, bucketName))
								_, err := fcClient.CreateTrigger(createFunctionResponse.Body.FunctionName, createTriggerRequest)
								if err != nil {
									if strings.Contains(err.Error(), "Cannot specify overlapping prefix and suffix with same event type.") {
										logger.Println.Error(fmt.Sprintf("在 %v 存储桶下已经存在相同前缀和后缀的触发器，此触发器无法创建。", bucketName))
										logger.Println.Info(fmt.Sprintf(`
已经创建的资源如下：
角色：%v
策略：%v
函数：%v
`, roleName, policyName, functionName))
									} else {
										logger.Println.Error("创建触发器时报错，详细信息如下：")
										logger.Println.Error(err.Error())
										logger.Println.Info(fmt.Sprintf(`
已经创建的资源如下：
角色：%v
策略：%v
函数：%v
`, roleName, policyName, functionName))
									}
								} else {
									logger.Println.Warn(fmt.Sprintf("限制存储桶只允许上传图片模块配置完成，现在向 %v 存储桶上传 jpg、jpeg、png、gif 之外格式的文件都会被自动删除。", bucketName))
									logger.Println.Warn(fmt.Sprintf(`已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
角色：%v
策略：%v
函数：%v
触发器：%v
`, roleName, policyName, functionName, triggerName))
								}
							}
						}
					}
				}
			}
		}
	}
}
