package utils

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func GetModuleName(m global.Module) string {
	return fmt.Sprintf("%d_%s_%s", m.ID, m.Provider.EnName, m.Name)
}

func GetModule(id, level int, provider, name, moduleProvider, introduce, moreInfo string, basicOptions []global.BasicOptions) global.Module {
	var (
		ids []int
	)
	sort.Slice(basicOptions, func(i, j int) bool {
		return basicOptions[i].Key < basicOptions[j].Key
	})

	// update global.BasicOptionsWithIds
	for _, BasicOptionsWithId := range global.BasicOptionsWithIds {
		ids = append(ids, BasicOptionsWithId.Id)
	}
	if !InInt(id, ids) {
		global.BasicOptionsWithIds = append(global.BasicOptionsWithIds, global.BasicOptionsWithId{
			Id:           id,
			BasicOptions: basicOptions,
		})
	}

	for n, BasicOptionsWithId := range global.BasicOptionsWithIds {
		if BasicOptionsWithId.Id == id {
			for _, fullBasicOption := range global.BasicOptionsFull {
				for m, basicOptionWithId := range BasicOptionsWithId.BasicOptions {
					if basicOptionWithId.Key == fullBasicOption.Key {
						global.BasicOptionsWithIds[n].BasicOptions[m].Value = global.GetBasicOptionValue(fullBasicOption.Key)
					}
				}
			}
		}
	}

	return global.Module{
		ID:             id,
		Provider:       getProvider(provider),
		Name:           name,
		ModuleProvider: moduleProvider,
		Introduce:      introduce,
		Level:          level,
		Info: generateDetailedInformation(id, level, getProvider(provider).Name, name, moduleProvider, introduce,
			GenerateTable(GetBasicOptions(global.GetBasicOptionsWithId(id))), moreInfo),
		BasicOptions: global.GetBasicOptionsWithId(id),
	}
}

func GetBasicOptions(basicOptions []global.BasicOptions) ([]string, [][]string, []int) {
	header := []string{"名称", "必选", "当前设置", "描述"}
	rows := [][]string{}
	for _, b := range basicOptions {
		if b.Required {
			rows = append(rows, []string{b.Key, global.True, b.Value, b.Introduce})
		} else {
			rows = append(rows, []string{b.Key, global.False, b.Value, b.Introduce})
		}
	}
	return header, rows, []int{15, 8, 40, 50}
}

func generateDetailedInformation(id, level int, provider, name, moduleProvider, introduce, basicOptions, moreInfo string) string {
	header := []string{"类型", "信息"}
	rows := [][]string{}
	rows = append(rows, []string{"ID", fmt.Sprintf("%d", id)})
	rows = append(rows, []string{"云提供商", provider})
	rows = append(rows, []string{"名称", name})
	rows = append(rows, []string{"推荐评级", strings.Repeat("★", level)})
	rows = append(rows, []string{"模块提供者", moduleProvider})
	rows = append(rows, []string{"模块简介", introduce})

	return fmt.Sprintf(`
 介绍：
%s

 操作：
%s
%s`, GenerateTable(header, rows, []int{15, 40}), basicOptions, moreInfo)
}

func getProvider(provider string) global.Provider {
	var p global.Provider
	switch provider {
	case global.Aliyun:
		p = global.Provider{"阿里云", global.Aliyun}
	case global.TencentCloud:
		p = global.Provider{"腾讯云", global.TencentCloud}
	case global.HuaweiCloud:
		p = global.Provider{"华为云", global.HuaweiCloud}
	case global.BaiduCloud:
		p = global.Provider{"百度云", global.BaiduCloud}
	}
	return p
}

// 文本处理

func Contains(arr []string, char string) bool {
	for _, c := range arr {
		if c == char {
			return true
		}
	}
	return false
}

func InInt(i int, int_array []int) bool {
	for _, element := range int_array {
		if i == element {
			return true
		}
	}
	return false
}

// 将空指针转为空字符串

func ConvertedNullPointer(p *string) string {
	if p == nil {
		return global.NULL
	} else {
		return *p
	}
}

// 如果为空则不输出

func PrintfNotNilString(p string, s *string) {
	if s != nil {
		fmt.Printf("%s%s\n", p, *s)
	}
}

func PrintfNotNilInt32(p string, s *int32) {
	if s != nil {
		fmt.Printf("%s%d\n", p, *s)
	}
}

func PrintfNotNilInt64(p string, s *int64) {
	if s != nil {
		fmt.Printf("%s%d\n", p, *s)
	}
}

func PrintfNotNilUInt64(p string, s *uint64) {
	if s != nil {
		fmt.Printf("%s%d\n", p, *s)
	}
}

func PrintfNotNilFloat64(p string, s *float64) {
	if s != nil {
		fmt.Printf("%s%f\n", p, *s)
	}
}

// 将 bytes 类型转为人类可阅读的形式

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// 让用户判断是否继续操作

func SurveyConfirm(message string) bool {
	var bool_ bool
	prompt := &survey.Confirm{
		Message: message,
		Default: true,
	}
	err := survey.AskOne(prompt, &bool_)
	if err != nil {
		logger.Println.Error(err.Error())
	}
	return bool_
}

// 生成随机名称

func GenerateRandomName(s string) string {
	return fmt.Sprintf("CloudSword%s%s", s, generateNumber(6))
}

// 生成指定位数的数字
func generateNumber(length int) string {
	const hexChars = "0123456789"
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(hexChars[rand.Intn(len(hexChars))])
	}

	return sb.String()
}

// 创建随机密码

func GenerateRandomPasswords() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "%@#$"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 16
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)
	return str
}
