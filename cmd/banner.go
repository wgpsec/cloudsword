package cmd

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils/global"
)

type EnvironmentVariables struct {
	Text         string
	Description  string
	BasicOptions string
}

var (
	banner = `
                                 (_)
                                 |_|         
                                 |_|         
                                 |_|         
  ▗▄▄▖▗▖    ▗▄▖ ▗▖ ▗▖▗▄▄▄    o=========o    ▗▄▄▖▗▖ ▗▖ ▗▄▖ ▗▄▄▖ ▗▄▄▄
 ▐▌   ▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌  █       | |       ▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌ ▐▌▐▌  █
 ▐▌   ▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌  █       | |        ▝▀▚▖▐▌ ▐▌▐▌ ▐▌▐▛▀▚▖▐▌  █
 ▝▚▄▄▖▐▙▄▄▖▝▚▄▞▘▝▚▄▞▘▐▙▄▄▀       | |       ▗▄▄▞▘▐▙█▟▌▝▚▄▞▘▐▌ ▐▌▐▙▄▄▀
                                 | |
                                 | |
                                 | |
                                 | |
                                 \ /
                      
                <-- 云鉴，让您的公有云环境更安全 -->
`
	projectInfo = fmt.Sprintf(`
                                %s
                             %s
           项目地址：https://github.com/wgpsec/cloudsword
           使用手册：https://wiki.teamssix.com/cloudsword
`, global.Version, global.UpdateDate)
	helpMessageEnvironmentVariables = []EnvironmentVariables{
		{
			Text:         global.CloudSwordAccessKeyID,
			BasicOptions: global.AKId,
			Description:  "访问凭证 ID",
		},
		{
			Text:         global.CloudSwordAccessKeySecret,
			BasicOptions: global.AKSecret,
			Description:  "访问凭证 Secret",
		},
		{
			Text:         global.CloudSwordSecurityToken,
			BasicOptions: global.AKToken,
			Description:  "可选，访问凭证的临时令牌部分",
		},
		//{
		//	Text:         global.CloudSwordProxy,
		//	BasicOptions: global.Proxy,
		//	Description:  "\t代理访问，支持 Socks5、HTTPS 和 HTTP",
		//},
		//{
		//	Text:         global.CloudSwordUserAgent,
		//	BasicOptions: global.UserAgent,
		//	Description:  "\t请求 UA 头",
		//},
		{
			Text:         global.CloudSwordDetail,
			BasicOptions: global.Detail,
			Description:  "\t详细内容输出（设置 no 或者 yes）",
		},
	}
)

func returnHelpInformation() string {
	var (
		level1               string
		level2               string
		example              string
		environmentVariables string
		helpInformation      string
		moreInformation      string
	)
	level1 = "\t全局命令\t描述\n\t--------\t----\n"
	level2 = "\t二级命令\t描述\n\t--------\t----\n"
	environmentVariables = "\t环境变量\t\t\t\t描述\n\t--------\t\t\t\t----\n"
	example = fmt.Sprintf(`
列出阿里云 OSS 对象存储桶：

%ssearch oss
%suse 1201_aliyun_oss_list_buckets
%sset ak_id xxxxxxx
%sset ak_secret xxxxxxx
%srun


创建腾讯云访问凭证蜜标：

%sexport CLOUD_SWORD_ACCESS_KEY_ID="xxxxxxx"
%sexport CLOUD_SWORD_ACCESS_KEY_SECRET="xxxxxxx"
%suse 2101_tencent_cloud_list_cloud_assets
%sset webhook https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx
%srun
`, global.Tab, global.Tab, global.Tab, global.Tab, global.Tab, global.Tab, global.Tab, global.Tab, global.Tab, global.Tab)
	moreInformation = fmt.Sprintf(`
当前版本：%s
发布日期：%s
项目地址：https://github.com/wgpsec/cloudsword
使用手册：https://wiki.teamssix.com/cloudsword
项目团队公众号：WgpSec 狼组安全团队
作者个人公众号：TeamsSix
==============================================
    建议关注公众号以获取云鉴相关的最新消息    
==============================================
`, global.Version, global.UpdateDate)
	for _, help := range helpMessageLevel1 {
		level1 += fmt.Sprintf("\t%s\t\t%s\n", help.Text, help.Description)
	}
	for _, help := range helpMessageLevel2 {
		level2 += fmt.Sprintf("\t%s\t\t%s\n", help.Text, help.Description)
	}
	for _, help := range helpMessageEnvironmentVariables {
		environmentVariables += fmt.Sprintf("\t%s\t\t%s\n", help.Text, help.Description)
	}

	helpInformation = fmt.Sprintf("%s\n\n全局命令\n========\n\n%s\n\n二级命令\n========\n\n%s\n\n环境变量\n========\n\n%s\n\n使用示例\n========\n%s", moreInformation,
		level1,
		level2,
		environmentVariables,
		example)
	return helpInformation
}
