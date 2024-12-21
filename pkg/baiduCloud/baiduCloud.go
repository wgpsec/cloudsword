package baiduCloud

import (
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
)

func Modules() []global.Module {
	var modules []global.Module
	//modules = append(modules, module1()...)
	modules = append(modules, module2()...)
	//modules = append(modules, module3()...)
	//modules = append(modules, module4()...)
	return modules
}

// OBS 模块

func module2() []global.Module {
	return []global.Module{
		utils.GetModule(
			4201,
			2,
			global.BaiduCloud,
			"bos_list_buckets",
			global.TeamsSix,
			"列出百度云 BOS 对象存储桶",
			`
所需权限：
1. bos:ListBuckets
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
	}
}
