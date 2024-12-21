package utils

import (
	"fmt"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"time"
)

func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func GetUTC8TimeType1(t string) (string, error) {
	inputTime, err := time.Parse("2006-01-02T15:04Z", t)
	if err != nil {
		logger.Println.Error(fmt.Sprintf("解析 %v 时间格式时报错，详细信息如下：", t))
		logger.Println.Error(err.Error())
		return global.NULL, err
	} else {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logger.Println.Error("修改时间地区时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			utc8Time := inputTime.In(loc)
			return utc8Time.Format("2006-01-02 15:04:05"), err
		}
	}
}

func GetUTC8TimeType2(t string) (string, error) {
	inputTime, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		logger.Println.Error(fmt.Sprintf("解析 %v 时间格式时报错，详细信息如下：", t))
		logger.Println.Error(err.Error())
		return global.NULL, err
	} else {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logger.Println.Error("修改时间地区时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			utc8Time := inputTime.In(loc)
			return utc8Time.Format("2006-01-02 15:04:05"), err
		}
	}
}

func GetUTC8TimeType3(t string) (string, error) {
	inputTime, err := time.Parse("2006-01-02T15:04:05.999Z", t)
	if err != nil {
		logger.Println.Error(fmt.Sprintf("解析 %v 时间格式时报错，详细信息如下：", t))
		logger.Println.Error(err.Error())
		return global.NULL, err
	} else {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			logger.Println.Error("修改时间地区时报错，详细信息如下：")
			logger.Println.Error(err.Error())
			return global.NULL, err
		} else {
			utc8Time := inputTime.In(loc)
			return utc8Time.Format("2006-01-02 15:04:05"), err
		}
	}
}
