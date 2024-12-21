package huaweiCloud

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
			3201,
			2,
			global.HuaweiCloud,
			"obs_list_buckets",
			global.TeamsSix,
			"列出华为云 OBS 对象存储桶",
			`
所需权限：
1. obs:ListBuckets
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
	}
}
