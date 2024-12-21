package tencentCloud

import (
	"fmt"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/wgpsec/cloudsword/utils/logger"
)

// 列出 CAM 用户

func getCAMUsers() ([]*cam.SubAccountInfo, error) {
	client, err := CAMClient()
	if err == nil {
		requests := cam.NewListUsersRequest()
		response, err := client.ListUsers(requests)
		if err != nil {
			logger.Println.Error("获取 CAM 用户时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			return nil, err
		} else {
			return response.Response.Data, err
		}
	} else {
		return nil, err
	}
}

// 列出 CAM 角色

func getCAMRoles() ([]*cam.RoleInfo, error) {
	client, err := CAMClient()
	if err == nil {
		var (
			rp    uint64 = 200
			page  uint64 = 1
			roles []*cam.RoleInfo
		)
		requests := cam.NewDescribeRoleListRequest()
		requests.Page = common.Uint64Ptr(page)
		requests.Rp = common.Uint64Ptr(rp)
		for {
			response, err := client.DescribeRoleList(requests)
			if err != nil {
				logger.Println.Error("获取 CAM 角色时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				break
			} else {
				roles = append(roles, response.Response.List...)
				if len(response.Response.List) == int(rp) {
					page = page + 1
					requests.Page = &page
				} else {
					break
				}
			}
		}
		return roles, err
	} else {
		return nil, err
	}
}

// 创建 CAM 用户

func addUser(userName, remark string) error {
	client, err := CAMClient()
	if err == nil {
		request := cam.NewAddUserRequest()
		request.Name = common.StringPtr(userName)
		request.Remark = common.StringPtr(remark)
		_, err = client.AddUser(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("创建 %v 用户时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
		} else {
			logger.Println.Info(fmt.Sprintf("用户 %v 创建成功。", userName))
		}
	}
	return err
}

// 列出所有策略

func listPolicies() ([]*cam.StrategyInfo, error) {
	client, err := CAMClient()
	if err == nil {
		var (
			page     uint64 = 1
			rp       uint64 = 200
			policies []*cam.StrategyInfo
		)

		requests := cam.NewListPoliciesRequest()
		requests.Rp = common.Uint64Ptr(rp)
		requests.Page = common.Uint64Ptr(page)
		requests.Scope = common.StringPtr("All")
		for {
			response, err := client.ListPolicies(requests)
			if err != nil {
				logger.Println.Error("获取策略时报错，详细信息如下：")
				logger.Println.Error(err.Error())
				break
			} else {
				policies = append(policies, response.Response.List...)
				if len(response.Response.List) == int(rp) {
					page = page + 1
					requests.Page = &page
				} else {
					break
				}
			}
		}
		return policies, err
	}
	return nil, err
}

// 查询用户信息

func getUser(userName string) (*cam.GetUserResponseParams, error) {
	client, err := CAMClient()
	if err == nil {
		request := cam.NewGetUserRequest()
		request.Name = common.StringPtr(userName)
		response, err := client.GetUser(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("查询 %v 用户信息时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
			return nil, err
		} else {
			logger.Println.Info(fmt.Sprintf("查询 %v 用户信息成功。", userName))
			return response.Response, err
		}
	} else {
		return nil, err
	}
}

// 为用户添加策略

func attachUserPolicy(userName, policyName string, uin, policyId uint64) error {
	client, err := CAMClient()
	if err == nil {
		request := cam.NewAttachUserPolicyRequest()
		request.AttachUin = common.Uint64Ptr(uin)
		request.PolicyId = common.Uint64Ptr(policyId)
		_, err := client.AttachUserPolicy(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("为 %v 用户添加 %v 策略时报错，详细信息如下：", userName, policyName))
			logger.Println.Error(err.Error())
		} else {
			logger.Println.Info(fmt.Sprintf("为 %v 用户添加 %v 策略成功。", userName, policyName))
		}
	}
	return err
}

// 为用户启用控制台登录

func updateUserWithConsoleLogin(userName, password string) error {
	client, err := CAMClient()
	if err == nil {
		request := cam.NewUpdateUserRequest()
		request.Name = common.StringPtr(userName)
		request.ConsoleLogin = common.Uint64Ptr(1)
		request.Password = common.StringPtr(password)
		request.NeedResetPassword = common.Uint64Ptr(0)
		_, err := client.UpdateUser(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("为 %v 用户创建 Web 控制台登录配置时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
		}
	}
	return err
}

// 为用户创建访问凭证

func createAccessKey(userName string, uin uint64) (*cam.AccessKeyDetail, error) {
	client, err := CAMClient()
	if err == nil {
		request := cam.NewCreateAccessKeyRequest()
		request.TargetUin = common.Uint64Ptr(uin)
		response, err := client.CreateAccessKey(request)
		if err != nil {
			logger.Println.Error(fmt.Sprintf("为 %v 用户创建访问凭证时报错，详细信息如下：", userName))
			logger.Println.Error(err.Error())
			return nil, err
		} else {
			return response.Response.AccessKey, err
		}
	} else {
		return nil, err
	}
}

// 获取用户 AppId

func getUserAppId() (*cam.GetUserAppIdResponseParams, error) {
	client, err := CAMClient()
	if err == nil {
		request := cam.NewGetUserAppIdRequest()
		response, err := client.GetUserAppId(request)
		if err != nil {
			logger.Println.Error("获取用户 AppId 时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			return nil, err
		} else {
			return response.Response, err
		}
	} else {
		return nil, err
	}
}
