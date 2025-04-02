package qiniuCloud

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

// Kodo 模块

func module2() []global.Module {
	return []global.Module{
		utils.GetModule(
			5201,
			2,
			global.QiniuCloud,
			"kodo_list_buckets",
			"ZhuriLab@moresec",
			"列出七牛云 Kodo 对象存储桶",
			`
所需权限：
无（七牛暂无权限管理，AK/SK 默认开放所有权限）
`,
			append(global.BasicOptionsDefault, global.BasicOptionDetail),
		),
	}
}
