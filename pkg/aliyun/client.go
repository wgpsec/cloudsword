package aliyun

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	domain "github.com/alibabacloud-go/domain-20180129/v4/client"
	fc "github.com/alibabacloud-go/fc-20230330/v4/client"
	ram "github.com/alibabacloud-go/ram-20150501/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	ossCredentials "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/transport"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"net/http"
)

// OSS

func OSSClient(region string) *oss.Client {
	var (
		//proxy        string
		cfg          *oss.Config
		customClient *http.Client
		transConfig  transport.Config
		transports   []func(*http.Transport)
	)

	cfg = oss.LoadDefaultConfig()

	cfg.WithCredentialsProvider(
		ossCredentials.NewStaticCredentialsProvider(
			global.GetBasicOptionValue(global.AKId),
			global.GetBasicOptionValue(global.AKSecret),
			global.GetBasicOptionValue(global.AKToken),
		))
	cfg.WithRegion(region)

	// 自定义 HTTP 客户端
	//proxy = global.GetBasicOptionValue(global.Proxy)
	//if proxy != global.NULL {
	//	transports = append(transports, utils.HttpTransport())
	//}
	customClient = transport.NewHttpClient(&transConfig, transports...)

	// 自定义 User Agent
	//customClient.Transport = &utils.UserAgentTransport{
	//	Base:      customClient.Transport,
	//	UserAgent: global.GetBasicOptionValue(global.UserAgent),
	//}

	cfg.WithHttpClient(customClient)
	client := oss.NewClient(cfg)
	return client
}

// ECS

func ECSClient(region string) (*ecs.Client, error) {
	var (
		ecsClient *ecs.Client
		err       error
	)
	ecsConfig := sdk.NewConfig()
	if global.GetBasicOptionValue(global.AKToken) == global.NULL {
		credential := credentials.NewStsTokenCredential(
			global.GetBasicOptionValue(global.AKId),
			global.GetBasicOptionValue(global.AKSecret),
			global.GetBasicOptionValue(global.AKToken),
		)
		ecsClient, err = ecs.NewClientWithOptions(region, ecsConfig, credential)
		if err != nil {
			logger.Println.Error("创建 ECS 客户端时报错，详细信息如下：")
			logger.Println.Error(err.Error())
		}
	} else {
		credential := credentials.NewAccessKeyCredential(
			global.GetBasicOptionValue(global.AKId),
			global.GetBasicOptionValue(global.AKSecret),
		)
		ecsClient, err = ecs.NewClientWithOptions(region, ecsConfig, credential)
		if err != nil {
			logger.Println.Error("创建 ECS 客户端时报错，详细信息如下：")
			logger.Println.Error(err.Error())
		}
	}
	return ecsClient, err
}

// RAM

func RAMClient() (*ram.Client, error) {
	credential := &openapi.Config{
		AccessKeyId:     tea.String(global.GetBasicOptionValue(global.AKId)),
		AccessKeySecret: tea.String(global.GetBasicOptionValue(global.AKSecret)),
		SecurityToken:   tea.String(global.GetBasicOptionValue(global.AKToken))}
	ramClient, err := ram.NewClient(credential)
	if err != nil {
		logger.Println.Error("创建 RAM 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
		return nil, err
	} else {
		return ramClient, err
	}
}

// Domains

func DomainClient() (*domain.Client, error) {
	credential := &openapi.Config{
		AccessKeyId:     tea.String(global.GetBasicOptionValue(global.AKId)),
		AccessKeySecret: tea.String(global.GetBasicOptionValue(global.AKSecret)),
		SecurityToken:   tea.String(global.GetBasicOptionValue(global.AKToken))}
	domainClient, err := domain.NewClient(credential)
	if err != nil {
		logger.Println.Error("创建 Domain 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
		return nil, err
	} else {
		return domainClient, err
	}
}

// FC

func FCClient(accountId, region string) (*fc.Client, error) {
	credential := &openapi.Config{
		AccessKeyId:     tea.String(global.GetBasicOptionValue(global.AKId)),
		AccessKeySecret: tea.String(global.GetBasicOptionValue(global.AKSecret)),
		SecurityToken:   tea.String(global.GetBasicOptionValue(global.AKToken))}
	credential.Endpoint = tea.String(fmt.Sprintf("%s.%s.fc.aliyuncs.com", accountId, region))
	fcClient, err := fc.NewClient(credential)
	if err != nil {
		logger.Println.Error("创建 FC 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
		return nil, err
	} else {
		return fcClient, err
	}
}
