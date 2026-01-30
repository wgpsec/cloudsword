package global

import (
	"fmt"
	"strconv"
)

// Version information - these variables can be set at build time using ldflags
var (
	Version    = "dev"      // Will be set by ldflags during build
	UpdateDate = "unknown"  // Will be set by ldflags during build
)

const (

	// 版本信息
	Team       = "WgpSec"
	Name       = "CloudSword"

	// 云提供商

	Aliyun       = "aliyun"
	TencentCloud = "tencent_cloud"
	HuaweiCloud  = "huawei_cloud"
	BaiduCloud   = "baidu_cloud"
	QiniuCloud   = "qiniu_cloud"

	// 基础选项

	AKId     = "ak_id"
	AKSecret = "ak_secret"
	AKToken  = "ak_token"
	//UserAgent = "user_agent"
	//Proxy     = "Proxy"
	Detail         = "detail"
	Region         = "region"
	ResultFilePath = "result_file_path"
	MaxResults     = "max_results"
	QueryValue     = "query_value"
	BucketName     = "bucket_name"
	UserName       = "user_name"
	Description    = "description"
	PolicyName     = "policy_name"
	Webhook        = "webhook"

	// 环境变量

	CloudSwordAccessKeyID     = "CLOUD_SWORD_ACCESS_KEY_ID"
	CloudSwordAccessKeySecret = "CLOUD_SWORD_ACCESS_KEY_SECRET"
	CloudSwordSecurityToken   = "CLOUD_SWORD_SECURITY_TOKEN"
	//CloudSwordUserAgent       = "CLOUD_SWORD_USER_AGENT"
	//CloudSwordProxy           = "CLOUD_SWORD_PROXY"
	CloudSwordDetail = "CLOUD_SWORD_DETAIL"

	// 布尔值

	True  = "true"
	False = "false"

	// 默认区域

	AliyunDefaultRegion       = "cn-hangzhou"
	TencentCloudDefaultRegion = "ap-guangzhou"
	HuaweiCloudDefaultRegion  = "cn-north-4"

	// 其他

	DefaultPrefix = Name + " > "
	TeamsSix      = "TeamsSix"
	Tab           = "\t"
	NULL          = ""
	All           = "all"
)

type Module struct {
	ID             int
	Provider       Provider
	Name           string
	ModuleProvider string
	Introduce      string
	Desc           string
	Level          int
	Info           string
	BasicOptions   []BasicOptions
}

type Provider struct {
	Name   string
	EnName string
}

type BasicOptions struct {
	Key       string
	Value     string
	Required  bool
	Introduce string
}

type BasicOptionsWithId struct {
	Id           int
	BasicOptions []BasicOptions
}

var BasicOptionsWithIds []BasicOptionsWithId

// BasicOptions

var BasicOptionsFull []BasicOptions

var BasicOptionsDefault = []BasicOptions{
	{
		Key:       AKId,
		Value:     "",
		Required:  true,
		Introduce: "访问凭证 ID",
	},
	{
		Key:       AKSecret,
		Value:     "",
		Required:  true,
		Introduce: "访问凭证 Secret",
	},
	{
		Key:       AKToken,
		Value:     "",
		Required:  false,
		Introduce: "可选，访问凭证的临时令牌部分",
	},
	//{
	//	Key:       Proxy,
	//	Value:     "",
	//	Required:  false,
	//	Introduce: "可选，代理访问，支持 Socks5、HTTPS 和 HTTP",
	//},
	//{
	//	Key:       UserAgent,
	//	Value:     fmt.Sprintf("%s/%s (%s; %s; %s)", Name, Version, runtime.Version(), runtime.GOOS, runtime.GOARCH),
	//	Required:  true,
	//	Introduce: "可选，请求 UA 头",
	//},
}

var BasicOptionDetail = BasicOptions{
	Key:       Detail,
	Value:     False,
	Required:  true,
	Introduce: "设置详细输出模式（true 或 false）",
}

var BasicOptionRegion = BasicOptions{
	Key:       Region,
	Value:     All,
	Required:  true,
	Introduce: "设置要列出的区域",
}

var BasicOptionResultFilePath = BasicOptions{
	Key:       ResultFilePath,
	Value:     NULL,
	Required:  false,
	Introduce: "设置结果导出路径，默认导出格式为 xlsx 格式",
}

var BasicOptionMaxResults = BasicOptions{
	Key:       MaxResults,
	Value:     strconv.Itoa(50),
	Required:  true,
	Introduce: "设置最大结果数量输出限制，-1 代表不限制",
}

var BasicOptionQueryValue = BasicOptions{
	Key:       QueryValue,
	Value:     NULL,
	Required:  true,
	Introduce: "设置要查询的内容",
}

var BasicOptionBucketName = BasicOptions{
	Key:       BucketName,
	Value:     NULL,
	Required:  true,
	Introduce: "设置 Bucket 名称",
}

var BasicOptionUserName = BasicOptions{
	Key:       UserName,
	Value:     NULL,
	Required:  true,
	Introduce: "设置用户名称",
}

var BasicOptionDescription = BasicOptions{
	Key:       Description,
	Value:     fmt.Sprintf("此资源由%v创建", Name),
	Required:  false,
	Introduce: "设置资源描述",
}

var BasicOptionPolicyName = BasicOptions{
	Key:       PolicyName,
	Value:     NULL,
	Required:  true,
	Introduce: "设置策略名称",
}

var BasicOptionWebhook = BasicOptions{
	Key:       Webhook,
	Value:     NULL,
	Required:  true,
	Introduce: "设置通知 Webhook 地址。",
}

//var MoreInfo = `
//Proxy 代理格式配置示例：
//
// http://127.0.0.1:8080
// https://127.0.0.1:8080
// socks5://127.0.0.1:1080
// socks5://user:password@127.0.0.1:1080
//`

func init() {
	BasicOptionsFull = append(BasicOptionsFull, BasicOptionsDefault...)
	BasicOptionsFull = append(
		BasicOptionsFull,
		BasicOptionDetail,
		BasicOptionRegion,
		BasicOptionQueryValue,
		BasicOptionMaxResults,
		BasicOptionResultFilePath,
		BasicOptionBucketName,
		BasicOptionUserName,
		BasicOptionDescription,
		BasicOptionPolicyName,
		BasicOptionWebhook,
	)
}

func GetBasicOptionsWithId(id int) []BasicOptions {
	var basicOptions []BasicOptions
	for _, option := range BasicOptionsWithIds {
		if option.Id == id {
			basicOptions = option.BasicOptions
		}
	}
	return basicOptions
}

func GetBasicOptionValue(key string) string {
	var value string
	for _, option := range BasicOptionsFull {
		if option.Key == key {
			value = option.Value
		}
	}
	return value
}

func UpdateBasicOptionValue(key string, value string) {
	for k, option := range BasicOptionsFull {
		if option.Key == key {
			BasicOptionsFull[k].Value = value
		}
	}
}
