package cmd

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/wgpsec/cloudsword/pkg/aliyun"
	"github.com/wgpsec/cloudsword/pkg/baiduCloud"
	"github.com/wgpsec/cloudsword/pkg/huaweiCloud"
	"github.com/wgpsec/cloudsword/pkg/qiniuCloud"
	"github.com/wgpsec/cloudsword/pkg/tencentCloud"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"os"
	"strconv"
	"strings"
)

// init

func init() {
	for _, Environment := range helpMessageEnvironmentVariables {
		value := os.Getenv(Environment.Text)
		if value != "" {
			for k, b := range global.BasicOptionsFull {
				if Environment.BasicOptions == b.Key {
					global.BasicOptionsFull[k].Value = value
				}
			}
		}
	}
}

func Modules() []global.Module {
	modules := []global.Module{}
	modules = append(modules, aliyun.Modules()...)
	modules = append(modules, tencentCloud.Modules()...)
	modules = append(modules, huaweiCloud.Modules()...)
	modules = append(modules, baiduCloud.Modules()...)
	modules = append(modules, qiniuCloud.Modules()...)
	return modules
}

// info

func getModuleByID(id int) global.Module {
	var mo global.Module
	for _, m := range Modules() {
		if m.ID == id {
			mo = m
		}
	}
	return mo
}

// list

func printModule(module []global.Module) {
	header := []string{"ID", "云提供商", "推荐评级", "名称", "介绍"}
	rows := [][]string{}
	for _, m := range module {
		rows = append(rows, []string{
			strconv.Itoa(m.ID),
			m.Provider.Name,
			strings.Repeat("★", m.Level),
			m.Name,
			m.Introduce,
		})
	}
	fmt.Println(utils.GenerateTable(header, rows, []int{4, 8, 10, 30, 50}))
}

// run

func runModule(m global.Module) {
	var bool_ bool = true
	for _, n := range m.BasicOptions {
		if n.Required && n.Value == global.NULL {
			logger.Println.Warn("未配置必选参数：" + n.Key)
			bool_ = false
		}
	}
	if bool_ {
		logger.Println.Info(fmt.Sprintf("正在运行 %s 模块。", utils.GetModuleName(m)))
		switch m.ID {
		// 阿里云
		case 1101:
			aliyun.ListCloudAssets()
		case 1201:
			aliyun.OSSListBuckets()
		case 1202:
			aliyun.OSSSearchObjects()
		case 1203:
			aliyun.OSSLimitBucketOnlyUploadImages()
		case 1301:
			aliyun.ECSListInstances()
		case 1401:
			aliyun.RAMListUsers()
		case 1402:
			aliyun.RAMListRoles()
		case 1403:
			aliyun.RAMCreateUser()
		case 1404:
			aliyun.RAMAttachPolicyToUser()
		case 1405:
			aliyun.RAMCreateLoginProfile()
		case 1406:
			aliyun.RAMCreateAccessKey()
		case 1501:
			aliyun.DomainListDomains()

		// 腾讯云
		case 2101:
			tencentCloud.ListCloudAssets()
		case 2102:
			tencentCloud.CreateHoneyToken()
		case 2201:
			tencentCloud.COSListBuckets()
		case 2301:
			tencentCloud.CVMListInstances()
		case 2302:
			tencentCloud.LHListInstances()
		case 2401:
			tencentCloud.CAMListUsers()
		case 2402:
			tencentCloud.CAMListRoles()
		case 2403:
			tencentCloud.CAMCreateUser()
		case 2404:
			tencentCloud.CAMAttachPolicyToUser()
		case 2405:
			tencentCloud.CAMCreateLoginProfile()
		case 2406:
			tencentCloud.CAMCreateAccessKey()
		// 华为云
		case 3201:
			huaweiCloud.OBSListBuckets()
		// 百度云
		case 4201:
			baiduCloud.BOSListBuckets()
		//七牛云
		case 5201:
			qiniuCloud.KodoListBuckets()
		}
	}
}

// search

func fuzzySearchModule(str string) {
	var result []global.Module
	for _, m := range Modules() {
		if strings.Contains(strconv.Itoa(m.ID), str) || strings.Contains(m.Provider.Name, str) || strings.Contains(m.Provider.EnName,
			str) || strings.Contains(m.Name, str) || strings.Contains(m.Introduce, str) {
			result = append(result, m)
		}
	}
	printModule(result)
}

// set

func returnModuleBasicOptionWithSet(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, m := range getModuleByID(useID).BasicOptions {
		s = append(s, prompt.Suggest{Text: m.Key, Description: m.Introduce})
	}
	input := strings.TrimPrefix(d.TextBeforeCursor(), "set ")
	input = strings.TrimSpace(input)
	return filterContains(s, input, true)
}

func returnModuleBasicOptionWithS(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, m := range getModuleByID(useID).BasicOptions {
		s = append(s, prompt.Suggest{Text: m.Key, Description: m.Introduce})
	}
	input := strings.TrimPrefix(d.TextBeforeCursor(), "s ")
	input = strings.TrimSpace(input)
	return filterContains(s, input, true)
}

// unset

func returnModuleBasicOptionByUnset(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, m := range getModuleByID(useID).BasicOptions {
		s = append(s, prompt.Suggest{Text: m.Key, Description: m.Introduce})
	}
	input := strings.TrimPrefix(d.TextBeforeCursor(), "unset ")
	input = strings.TrimSpace(input)
	return filterContains(s, input, true)
}

// use

func preciseSearchModule(str string) []global.Module {
	var result []global.Module
	for _, m := range Modules() {
		if utils.GetModuleName(m) == str {
			result = append(result, m)
		}
	}
	return result
}

func returnAllModulesWithUse(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, m := range Modules() {
		s = append(s, prompt.Suggest{Text: utils.GetModuleName(m), Description: m.Introduce})
	}
	input := strings.TrimPrefix(d.TextBeforeCursor(), "use ")
	input = strings.TrimSpace(input)
	return filterContains(s, input, true)
}

func returnAllModulesWithU(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	for _, m := range Modules() {
		s = append(s, prompt.Suggest{Text: utils.GetModuleName(m), Description: m.Introduce})
	}
	input := strings.TrimPrefix(d.TextBeforeCursor(), "u ")
	input = strings.TrimSpace(input)
	return filterContains(s, input, true)
}

func filterContains(suggestions []prompt.Suggest, word string, ignoreCase bool) []prompt.Suggest {
	if len(word) == 0 {
		return suggestions
	}
	filtered := make([]prompt.Suggest, 0, len(suggestions))
	for _, s := range suggestions {
		text := s.Text
		if ignoreCase {
			text = strings.ToLower(text)
			word = strings.ToLower(word)
		}
		if strings.Contains(text, word) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
