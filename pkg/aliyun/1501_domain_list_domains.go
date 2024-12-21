package aliyun

import (
	"encoding/json"
	"fmt"
	domain "github.com/alibabacloud-go/domain-20180129/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// domain_list_domains

func DomainListDomains() {
	var (
		page    int32 = 1
		domains []*domain.QueryDomainListResponseBodyDataDomain
	)
	queryDomainListRequest := &domain.QueryDomainListRequest{}
	queryDomainListRequest.PageNum = tea.Int32(page)
	queryDomainListRequest.PageSize = tea.Int32(10)
	runtime := &util.RuntimeOptions{}
	for {
		client, err := DomainClient()
		if err == nil {
			response, err := client.QueryDomainListWithOptions(queryDomainListRequest, runtime)
			if err != nil {
				logger.Println.Error("列出域名时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				break
			} else {
				domains = append(domains, response.Body.Data.Domain...)
				if *response.Body.NextPage {
					page = page + 1
					queryDomainListRequest.PageNum = tea.Int32(page)
				} else {
					break
				}
			}
		}
	}

	if len(domains) == 0 {
		logger.Println.Info("未找到域名。\n")
	} else {
		logger.Println.Info("找到以下域名：")
		for _, domain_ := range domains {
			if global.GetBasicOptionValue(global.Detail) == global.False {
				fmt.Println(*domain_.DomainName)
			} else {
				utils.PrintfNotNilString("\n域名：", domain_.DomainName)
				utils.PrintfNotNilString("公司名称：", domain_.Ccompany)
				fmt.Printf("状态：%v\n", domainStatus(*domain_.DomainStatus))
				fmt.Printf("是否过期：%v\n", expirationDateStatus(*domain_.ExpirationDateStatus))
				utils.PrintfNotNilString("实名认证状态：", domain_.DomainAuditStatus)
				utils.PrintfNotNilString("分组编号：", domain_.DomainGroupId)
				utils.PrintfNotNilString("分组名称：", domain_.DomainGroupId)
				utils.PrintfNotNilString("资源组 ID：", domain_.ResourceGroupId)

				jsonBytes, err := json.Marshal(&domain_.Tag.Tag)
				if err != nil {
					logger.Println.Error("解析域名标签信息时报错，详细信息如下：")
					logger.Println.Error(err.Error())
				} else {
					fmt.Printf("标签：%s\n", jsonBytes)
				}

				utils.PrintfNotNilString("备注：", domain_.Remark)
				utils.PrintfNotNilString("注册时间：", domain_.RegistrationDate)
				utils.PrintfNotNilString("到期时间：", domain_.ExpirationDate)
				utils.PrintfNotNilInt32("距离到期天数：", domain_.ExpirationCurrDateDiff)
			}
		}
		fmt.Println()
	}
}

func expirationDateStatus(s string) string {
	var status string
	if s == "1" {
		status = "未过期"
	} else if s == "2" {
		status = "已过期"
	}
	return status
}

func domainStatus(s string) string {
	var status string
	if s == "1" {
		status = "急需续费"
	} else if s == "2" {
		status = "急需赎回"
	} else if s == "3" {
		status = "正常"
	}
	return status
}
