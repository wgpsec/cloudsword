package tencentCloud

import (
	"fmt"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strings"
)

func CreateHoneyToken() {
	// 判断用户输入的 Webhook 类型
	var noticeType string
	webhook := global.GetBasicOptionValue(global.Webhook)
	switch {
	case strings.Contains(webhook, "weixin.qq.com"):
		noticeType = "WeCom"
	case strings.Contains(webhook, "dingtalk.com"):
		noticeType = "DingTalk"
	case strings.Contains(webhook, "feishu.cn"):
		noticeType = "Lark"
	default:
		noticeType = "Http"
	}
	// 1. 创建用户
	userName := utils.GenerateRandomName("User")
	err := addUser(userName, "此用户由云鉴 2102_create_honey_token 模块自动创建，此用户用于生成访问凭证蜜标使用。")
	if err == nil {
		// 2. 创建访问凭证
		getUserResponseParams, err := getUser(userName)
		if err == nil {
			accessKeyDetail, err := createAccessKey(userName, *getUserResponseParams.Uin)
			if err != nil {
				fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
`, userName)
			} else {
				accessKeyId := *accessKeyDetail.AccessKeyId
				accessKeySecret := *accessKeyDetail.SecretAccessKey
				logger.Println.Info(fmt.Sprintf("为 %v 用户创建访问凭证成功。", userName))
				// 3. 创建日志集
				clsClient, err := CLSClient(global.TencentCloudDefaultRegion)
				if err == nil {
					logSetName := utils.GenerateRandomName("LogSet")
					createLogSetRequest := cls.NewCreateLogsetRequest()
					createLogSetRequest.LogsetName = common.StringPtr(logSetName)
					createLogSetResponse, err := clsClient.CreateLogset(createLogSetRequest)
					if err != nil {
						logger.Println.Error("创建日志集时报错，详细信息如下：")
						logger.Println.Error(err.Error())
						fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
`, userName)
					} else {
						logger.Println.Info(fmt.Sprintf("日志集 %v 创建成功。", logSetName))
						logSetId := *createLogSetResponse.Response.LogsetId
						// 4. 创建日志主题
						createTopicRequest := cls.NewCreateTopicRequest()
						topicName := utils.GenerateRandomName("Topic")
						createTopicRequest.TopicName = common.StringPtr(topicName)
						createTopicRequest.LogsetId = common.StringPtr(logSetId)
						createTopicRequest.Describes = common.StringPtr("此日志主题由云鉴 2102_create_honey_token 模块自动创建，此日志主题用于接收跟踪集传过来的数据。")
						createTopicResponse, err := clsClient.CreateTopic(createTopicRequest)
						if err != nil {
							logger.Println.Error("创建日志主题时报错，详细信息如下：")
							logger.Println.Error(err.Error())
							fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
日志服务-日志集：%s
`, userName, logSetName)
						} else {
							logger.Println.Info(fmt.Sprintf("日志主题 %v 创建成功。", topicName))
							TopicId := *createTopicResponse.Response.TopicId
							// 5. 创建索引
							createIndex := cls.NewCreateIndexRequest()
							createIndex.TopicId = common.StringPtr(TopicId)
							createIndex.Status = common.BoolPtr(true)
							createIndex.MetadataFlag = common.Uint64Ptr(1)
							createIndex.IncludeInternalFields = common.BoolPtr(true)
							createIndex.Rule = &cls.RuleInfo{
								DynamicIndex: &cls.DynamicIndex{
									Status: common.BoolPtr(true),
								},
								FullText: &cls.FullTextInfo{
									CaseSensitive: common.BoolPtr(false),
									ContainZH:     common.BoolPtr(true),
									Tokenizer:     common.StringPtr(`@&?|#()='\",;:<>[]{}/ \n\t\r\\`),
								},
								KeyValue: &cls.RuleKeyValueInfo{},
							}
							_, err := clsClient.CreateIndex(createIndex)
							if err != nil {
								logger.Println.Error("创建索引时报错，详细信息如下：")
								logger.Println.Error(err.Error())
							} else {
								logger.Println.Info(fmt.Sprintf("日志主题 %v 索引创建成功。", topicName))
								// 6. 创建跟踪集
								cloudAuditClient, err := CloudAuditClient()
								if err == nil {
									auditTrackName := utils.GenerateRandomName("AuditTrack")
									createAuditTrackRequest := cloudaudit.NewCreateAuditTrackRequest()
									createAuditTrackRequest.Name = common.StringPtr(auditTrackName)
									createAuditTrackRequest.Status = common.Uint64Ptr(1)
									createAuditTrackRequest.Storage = &cloudaudit.Storage{
										StorageType:   common.StringPtr("cls"),
										StorageRegion: common.StringPtr(global.TencentCloudDefaultRegion),
										StorageName:   common.StringPtr(TopicId),
										StoragePrefix: common.StringPtr("/"),
									}
									createAuditTrackRequest.ActionType = common.StringPtr("*")
									createAuditTrackRequest.ResourceType = common.StringPtr("*")
									createAuditTrackRequest.EventNames = common.StringPtrs([]string{"*"})
									_, err := cloudAuditClient.CreateAuditTrack(createAuditTrackRequest)
									if err != nil {
										logger.Println.Error("创建跟踪集时报错，详细信息如下：")
										logger.Println.Error(err.Error())
										fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
日志服务-日志集：%s
日志服务-日志主题：%s
`, userName, logSetName, topicName)
									} else {
										logger.Println.Info(fmt.Sprintf("跟踪集 %v 创建成功。", auditTrackName))
										// 7. 创建通知内容模版
										noticeContentName := utils.GenerateRandomName("NoticeContent")
										createNoticeContentRequest := cls.NewCreateNoticeContentRequest()
										createNoticeContentRequest.Name = common.StringPtr(noticeContentName)
										createNoticeContentRequest.Type = common.Uint64Ptr(0)
										createNoticeContentRequest.NoticeContents = []*cls.NoticeContent{
											{
												Type: common.StringPtr(noticeType),
												TriggerContent: &cls.NoticeContentInfo{
													Title: common.StringPtr("【紧急告警】蜜标凭证被调用，发现入侵行为"),
													Content: common.StringPtr(`
告警策略：{{.Alarm}}
告警级别：{{.Level_zh}}
攻击 IP：{{.QueryLog[0][0].content.sourceIPAddress}}
攻击区域：{{.QueryLog[0][0].content.eventRegion}}
攻击服务：{{.QueryLog[0][0].content.resourceType}}
攻击操作：{{.QueryLog[0][0].content.eventName}}
攻击结果：{{.QueryLog[0][0].content.errorMessage}}
攻击请求头：{{.QueryLog[0][0].content.userAgent}}
触发时间：{{.StartTime}}

详细报告：[{{.DetailUrl}}]({{.DetailUrl}})
查询数据：[{{.QueryUrl}}]({{.QueryUrl}})
{{- if .CanSilent}}
屏蔽告警：[{{.SilentUrl}}]({{.SilentUrl}})
{{- end}}
`),
												},
												RecoveryContent: &cls.NoticeContentInfo{
													Title: common.StringPtr("攻击者已停止调用蜜标凭证"),
													Content: common.StringPtr(`
告警策略：{{.Alarm}}
告警级别：{{.Level_zh}}
触发时间：{{.StartTime}}
恢复时间：{{.NotifyTime}}

详细报告：[{{.DetailUrl}}]({{.DetailUrl}})
查询数据：[{{.QueryUrl}}]({{.QueryUrl}})
{{- if .CanSilent}}
屏蔽告警：[{{.SilentUrl}}]({{.SilentUrl}})
{{- end}}
`),
												},
											},
										}
										createNoticeContentResponse, err := clsClient.CreateNoticeContent(createNoticeContentRequest)
										if err != nil {
											logger.Println.Error("创建通知模版时报错，详细信息如下：")
											logger.Println.Error(err.Error())
											fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
操作审计-跟踪集：%s
日志服务-日志集：%s
日志服务-日志主题：%s
`, userName, auditTrackName, logSetName, topicName)
										} else {
											logger.Println.Info(fmt.Sprintf("通知模块 %v 创建成功。", noticeContentName))
											noticeContentId := *createNoticeContentResponse.Response.NoticeContentId
											// 8. 创建通知渠道组
											alarmNoticeName := utils.GenerateRandomName("AlarmNotice")
											createAlarmNoticeRequest := cls.NewCreateAlarmNoticeRequest()
											createAlarmNoticeRequest.Name = common.StringPtr(alarmNoticeName)
											createAlarmNoticeRequest.NoticeRules = []*cls.NoticeRule{
												{
													Rule: common.StringPtr("{\"Value\":\"AND\",\"Type\":\"Operation\",\"Children\":[{\"Type\":\"Condition\",\"Value\":\"Level\",\"Children\":[{\"Value\":\"In\",\"Type\":\"Compare\"},{\"Value\":\"[2]\",\"Type\":\"Value\"}]}]}"),
													WebCallbacks: []*cls.WebCallback{
														{
															CallbackType:    common.StringPtr(noticeType),
															Url:             common.StringPtr(webhook),
															NoticeContentId: common.StringPtr(noticeContentId),
														},
													},
												}}
											createAlarmNoticeResponse, err := clsClient.CreateAlarmNotice(createAlarmNoticeRequest)
											if err != nil {
												logger.Println.Error("创建通知渠道组时报错，详细信息如下：")
												logger.Println.Error(err.Error())
												fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
操作审计-跟踪集：%s
日志服务-日志集：%s
日志服务-日志主题：%s
日志服务-通知模版：%s
`, userName, auditTrackName, logSetName, topicName, noticeContentName)
											} else {
												alarmNoticeId := *createAlarmNoticeResponse.Response.AlarmNoticeId
												// 9. 创建告警策略
												alarmName := utils.GenerateRandomName("Alarm")
												createAlarmRequest := cls.NewCreateAlarmRequest()
												createAlarmRequest.Name = common.StringPtr(alarmName)
												createAlarmRequest.AlarmTargets = []*cls.AlarmTarget{
													{
														Query:           common.StringPtr(fmt.Sprintf("userIdentity.secretId:%s", accessKeyId)),
														Number:          common.Int64Ptr(1),
														StartTimeOffset: common.Int64Ptr(-10),
														EndTimeOffset:   common.Int64Ptr(0),
														LogsetId:        common.StringPtr(logSetId),
														TopicId:         common.StringPtr(TopicId),
														SyntaxRule:      common.Uint64Ptr(1),
													},
												}
												createAlarmRequest.MonitorTime = &cls.MonitorTime{
													Type: common.StringPtr("Period"),
													Time: common.Int64Ptr(5),
												}
												createAlarmRequest.TriggerCount = common.Int64Ptr(1)
												createAlarmRequest.AlarmPeriod = common.Int64Ptr(5)
												createAlarmRequest.AlarmNoticeIds = common.StringPtrs([]string{alarmNoticeId})
												createAlarmRequest.AlarmLevel = common.Uint64Ptr(2)
												createAlarmRequest.Condition = common.StringPtr("[$1.__QUERYCOUNT__]> 0")
												createAlarmRequest.Status = common.BoolPtr(true)
												createAlarmRequest.GroupTriggerStatus = common.BoolPtr(false)
												_, err := clsClient.CreateAlarm(createAlarmRequest)
												if err != nil {
													logger.Println.Error("创建告警策略时报错，详细信息如下：")
													logger.Println.Error(err.Error())
												} else {
													logger.Println.Info(fmt.Sprintf("告警策略 %v 创建成功。", alarmName))
													logger.Println.Info("云访问凭证蜜标创建成功，您的云访问凭证蜜标如下：")
													fmt.Printf(`
AccessKeyId: %v
AccessKeySecret: %v

注意事项：
1. 这个功能所返回的访问凭证是没有添加任何权限的，但是为了避免一些意外情况，建议在非生产环境下使用此功能。
2. 这个功能初始化需要 10 分钟左右的时间，当初始化完成后，从蜜标凭证被调用到接收到告警大概会经历 5 分钟左右的时间。
`, accessKeyId, accessKeySecret)
												}
												fmt.Printf(`
已经创建的资源如下，如果想取消使用这个模块，请手动删除以下资源：
访问管理-用户：%s
操作审计-跟踪集：%s
日志服务-日志集：%s
日志服务-日志主题：%s
日志服务-通知模版：%s
日志服务-通知渠道组：%s
日志服务-告警策略：%s
`, userName, auditTrackName, logSetName, topicName, noticeContentName, alarmNoticeName, alarmName)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
