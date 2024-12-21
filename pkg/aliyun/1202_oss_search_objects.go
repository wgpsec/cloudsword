package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"strconv"
	"time"
)

// oss_search_objects

func OSSSearchObjects() {
	// 1. 查询是否启用数据索引
	bucketName := global.GetBasicOptionValue(global.BucketName)
	region, err := getBucketRegion(bucketName)
	if err == nil {
		status, err := getMetaQueryStatus(bucketName, region)
		if err == nil {
			// 2. 如果没有就启用数据索引
			if status == "MetaQueryNotExist" {
				if utils.SurveyConfirm("查询 OSS 对象需要先创建数据索引，这需要一些时间，并且在查询过程中可能会产生一些费用，您要继续吗？") {
					err = openMetaQuery(bucketName, region)
					if err == nil {
						logger.Printf.Infof("正在创建数据索引")
						for {
							status, err = getMetaQueryStatus(bucketName, region)
							if err == nil {
								if status == "Ready" {
									fmt.Printf(".")
								} else if status == "Running" {
									fmt.Println()
									break
								}
								time.Sleep(10 * time.Second)
							} else {
								break
							}
						}
					}
				}
			}

			// 3. 数据索引创建完成，开始查询数据
			if status == "Running" {
				logger.Println.Info("数据索引已建立，正在查询数据...")
				var results []oss.MetaQueryFile
				doMetaQueryRequest := &oss.DoMetaQueryRequest{
					MetaQuery: &oss.MetaQuery{},
				}
				//doMetaQueryRequest.MetaQuery
				doMetaQueryRequest.Bucket = tea.String(bucketName)

				maxResults, err := strconv.ParseInt(global.GetBasicOptionValue(global.MaxResults), 10, 64)
				if err != nil {
					logger.Println.Error("解析最大结果数量值时报错，详细信息如下：")
					logger.Println.Error(err.Error())
				} else {
					if maxResults == -1 {
						doMetaQueryRequest.MetaQuery.MaxResults = tea.Int64(100)
					} else {
						doMetaQueryRequest.MetaQuery.MaxResults = tea.Int64(maxResults)
					}
					doMetaQueryRequest.MetaQuery.Query = tea.String(fmt.Sprintf("{\"Field\": \"Filename\",\"Value\": \"%v\",\"Operation\": \"match\"}",
						global.GetBasicOptionValue(global.QueryValue)))
					for {
						responses, err := OSSClient(region).DoMetaQuery(context.TODO(), doMetaQueryRequest)
						if err != nil {
							logger.Println.Error(fmt.Sprintf("查询 %v 存储桶对象时报错错，详细信息如下："), bucketName)
							logger.Println.Error(err.Error())
							break
						} else {
							results = append(results, responses.Files...)
							if responses.NextToken == nil {
								break
							} else if *responses.NextToken == global.NULL {
								break
							} else {
								if maxResults == -1 {
									doMetaQueryRequest.MetaQuery.NextToken = responses.NextToken
								} else if len(results) > int(maxResults) {
									break
								} else {
									doMetaQueryRequest.MetaQuery.NextToken = responses.NextToken
								}
							}
						}
					}

					if len(results) == 0 {
						logger.Println.Info("未查询到数据。")
					} else {
						logger.Println.Info(fmt.Sprintf("查询到 %v 条数据，结果如下：", len(results)))
						for _, file := range results {
							if global.GetBasicOptionValue(global.Detail) == global.False {
								fmt.Println(*file.Filename)
							} else {
								fmt.Printf("\n名称：%v\n", *file.Filename)
								fmt.Printf("大小：%v\n", utils.FormatBytes(*file.Size))

								jsonBytes, err := json.Marshal(file.OSSTagging)
								if err != nil {
									logger.Println.Error(fmt.Sprintf("解析 %v 对象标签信息时报错，详细信息如下：", *file.Filename))
									logger.Println.Error(err.Error())
								} else {
									fmt.Printf("标签：%v\n", string(jsonBytes))
								}

								jsonBytes, err = json.Marshal(file.OSSUserMeta)
								if err != nil {
									logger.Println.Error(fmt.Sprintf("解析 %v 对象自定义元数据列表时报错，详细信息如下：", *file.Filename))
									logger.Println.Error(err.Error())
								} else {
									fmt.Printf("自定义元数据列表：%v\n", string(jsonBytes))
								}
								logger.Println.Error(err.Error())

								fmt.Printf("访问地址：https://%v.oss-%v.aliyuncs.com/%v\n", bucketName, region, *file.Filename)
							}
						}
						fmt.Println()
					}
					// 4. 关闭数据索引
					if utils.SurveyConfirm("需要关闭数据索引吗？（如果您当前没有搜索需求，建议选择关闭）") {
						closeMetaQuery(bucketName, region)
					}
				}
			}
		}
	}
}
