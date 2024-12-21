package tencentCloud

import (
	"fmt"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	lh "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"net/http"
	"net/url"
)

// COS

func COSClient(bucketName, region, nextMarker string) *cos.Client {
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if bucketName == global.NULL || region == global.NULL {
		cosClient := cos.NewClient(nil, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:     AKId,
				SecretKey:    AKSecret,
				SessionToken: AKToken,
			},
		})
		return cosClient
	} else {
		if nextMarker == global.NULL {
			u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, region))
			b := &cos.BaseURL{BucketURL: u}
			cosClient := cos.NewClient(b, &http.Client{
				Transport: &cos.AuthorizationTransport{
					SecretID:     AKId,
					SecretKey:    AKSecret,
					SessionToken: AKToken,
				},
			})
			return cosClient
		} else {
			u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com/?marker=%s", bucketName, region, nextMarker))
			b := &cos.BaseURL{BucketURL: u}
			cosClient := cos.NewClient(b, &http.Client{
				Transport: &cos.AuthorizationTransport{
					SecretID:     AKId,
					SecretKey:    AKSecret,
					SessionToken: AKToken,
				},
			})
			return cosClient
		}
	}
}

// CVM

func CVMClient(region string) (*cvm.Client, error) {
	var (
		credential *common.Credential
		client     *cvm.Client
		err        error
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken == global.NULL {
		credential = common.NewCredential(AKId, AKSecret)
	} else {
		credential = common.NewTokenCredential(AKId, AKSecret, AKToken)
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	client, err = cvm.NewClient(credential, region, cpf)
	if err != nil {
		logger.Println.Error("创建 CVM 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return client, err
}

// LH

func LHClient(region string) (*lh.Client, error) {
	var (
		credential *common.Credential
		client     *lh.Client
		err        error
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken == global.NULL {
		credential = common.NewCredential(AKId, AKSecret)
	} else {
		credential = common.NewTokenCredential(AKId, AKSecret, AKToken)
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "lighthouse.tencentcloudapi.com"
	client, err = lh.NewClient(credential, region, cpf)
	if err != nil {
		logger.Println.Error("创建 LH 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return client, err
}

// CAM

func CAMClient() (*cam.Client, error) {
	var (
		credential *common.Credential
		client     *cam.Client
		err        error
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken == global.NULL {
		credential = common.NewCredential(AKId, AKSecret)
	} else {
		credential = common.NewTokenCredential(AKId, AKSecret, AKToken)
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cam.tencentcloudapi.com"
	client, err = cam.NewClient(credential, global.NULL, cpf)
	if err != nil {
		logger.Println.Error("创建 CAM 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return client, err
}

// CloudAudit

func CloudAuditClient() (*cloudaudit.Client, error) {
	var (
		credential *common.Credential
		client     *cloudaudit.Client
		err        error
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken == global.NULL {
		credential = common.NewCredential(AKId, AKSecret)
	} else {
		credential = common.NewTokenCredential(AKId, AKSecret, AKToken)
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cloudaudit.tencentcloudapi.com"
	client, err = cloudaudit.NewClient(credential, global.TencentCloudDefaultRegion, cpf)
	if err != nil {
		logger.Println.Error("创建 Cloud Audit 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return client, err
}

// CLS

func CLSClient(region string) (*cls.Client, error) {
	var (
		credential *common.Credential
		client     *cls.Client
		err        error
	)
	AKId := global.GetBasicOptionValue(global.AKId)
	AKSecret := global.GetBasicOptionValue(global.AKSecret)
	AKToken := global.GetBasicOptionValue(global.AKToken)
	if AKToken == global.NULL {
		credential = common.NewCredential(AKId, AKSecret)
	} else {
		credential = common.NewTokenCredential(AKId, AKSecret, AKToken)
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cls.tencentcloudapi.com"
	client, err = cls.NewClient(credential, region, cpf)
	if err != nil {
		logger.Println.Error("创建 CLS 客户端时报错，详细信息如下：")
		logger.Println.Error(err.Error())
	}
	return client, err
}
