package cmd

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"os"
	"strings"
)

type helpStruct struct {
	Text        string
	Description string
}

var (
	useID             int
	currentPrefix     = global.DefaultPrefix
	p                 *prompt.Prompt
	helpMessageLevel1 = []helpStruct{
		{Text: "help", Description: "查看帮助信息"},
		{Text: "list", Description: "列出模块"},
		{Text: "quit", Description: "退出程序"},
		{Text: "search", Description: "搜索模块"},
		{Text: "use", Description: "使用模块"},
	}
	helpMessageLevel2 = []helpStruct{
		{Text: "info", Description: "查看模块使用方法"},
		{Text: "run", Description: "运行模块"},
		{Text: "set", Description: "设置运行参数"},
		{Text: "unset", Description: "取消设置运行参数"},
	}
)

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	spaceCount := strings.Count(in.TextBeforeCursor(), " ")
	// use
	if strings.HasPrefix(in.TextBeforeCursor(), "use ") {
		if spaceCount < 2 {
			return returnAllModulesWithUse(in)
		} else {
			return nil
		}
	}
	// u
	if strings.HasPrefix(in.TextBeforeCursor(), "u ") {
		if spaceCount < 2 {
			return returnAllModulesWithU(in)
		} else {
			return nil
		}
	}

	// set
	if strings.HasPrefix(in.TextBeforeCursor(), "set ") {
		if spaceCount < 2 {
			return returnModuleBasicOptionWithSet(in)
		} else {
			return nil
		}
	}
	// s
	if strings.HasPrefix(in.TextBeforeCursor(), "s ") {
		if spaceCount < 2 {
			return returnModuleBasicOptionWithS(in)
		} else {
			return nil
		}
	}

	// unset
	if strings.HasPrefix(in.TextBeforeCursor(), "unset ") {
		if spaceCount < 2 {
			return returnModuleBasicOptionByUnset(in)
		} else {
			return nil
		}
	}
	// other
	if spaceCount < 1 {
		if useID == 0 {
			for _, h := range helpMessageLevel1 {
				s = append(s, prompt.Suggest{Text: h.Text, Description: h.Description})
			}
		} else {
			for _, h := range helpMessageLevel1 {
				s = append(s, prompt.Suggest{Text: h.Text, Description: h.Description})
			}
			for _, h := range helpMessageLevel2 {
				s = append(s, prompt.Suggest{Text: h.Text, Description: h.Description})
			}
		}
		return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
	}
	return nil
}

func changingPrefixLivePrefix() (string, bool) {
	return currentPrefix, true
}

func executor(in string) {
	in = strings.TrimSpace(in)
	var blocks []string
	blocks_ := strings.Split(in, " ")
	for _, block := range blocks_ {
		block = strings.Trim(block, " \"")
		if len(block) > 0 {
			blocks = append(blocks, block)
		}
	}
	if len(blocks) > 0 {
		switch blocks[0] {
		case "help", "h":
			fmt.Println(returnHelpInformation())
		case "info", "i":
			if useID == 0 {
				logger.Println.Warn("请使用 use 命令选择你要调用的模块。\n")
			} else {
				fmt.Println(getModuleByID(useID).Info)
			}
		case "list", "l":
			printModule(Modules())
		case "quit", "q":
			logger.Println.Warn("云鉴已退出。")
			os.Exit(0)
		case "run", "r":
			if useID == 0 {
				logger.Println.Warn("请使用 use 命令选择你要调用的模块。\n")
			} else if global.GetBasicOptionValue(global.AKId) == global.NULL {
				logger.Println.Warn("请配置云访问凭证后再运行。\n")
			} else {
				runModule(getModuleByID(useID))
			}
		case "search":
			if len(blocks) == 2 {
				fuzzySearchModule(blocks[1])
			} else {
				logger.Println.Warn("请检查你的命令。\n")
			}
		case "set", "s":
			if useID == 0 {
				logger.Println.Warn("请使用 use 命令选择你要调用的模块。\n")
			} else {
				if len(blocks) == 3 {
					for _, b := range global.GetBasicOptionsWithId(useID) {
						if blocks[1] == b.Key {
							if b.Key == global.Detail {
								value := strings.ToLower(blocks[2])
								if utils.Contains([]string{"0", "false", "f", "no", "n"}, value) {
									global.UpdateBasicOptionValue(b.Key, global.False)
									fmt.Println(b.Key + " ==> " + global.False)
								} else if utils.Contains([]string{"1", "true", "t", "yes", "y"}, value) {
									global.UpdateBasicOptionValue(b.Key, global.True)
									fmt.Println(b.Key + " ==> " + global.True)
								} else {
									logger.Println.Error(fmt.Sprintf("%s 的值类型错误，请设置成 True 或者 False。", b.Key))
								}
							} else {
								global.UpdateBasicOptionValue(b.Key, blocks[2])
								fmt.Println(b.Key + " ==> " + blocks[2])
							}
						}
					}
				} else {
					logger.Println.Warn("请检查你的命令，正确格式：set <key> <value>。\n")
				}
			}
		case "unset":
			if useID == 0 {
				logger.Println.Warn("请使用 use 命令选择你要调用的模块。\n")
			} else {
				if len(blocks) == 2 {
					for k, b := range global.GetBasicOptionsWithId(useID) {
						if blocks[1] == b.Key {
							global.BasicOptionsFull[k].Value = global.NULL
							fmt.Println(b.Key + " ==> " + global.NULL)
						}
					}
				} else {
					logger.Println.Warn("请检查你的命令。\n")
				}
			}
		case "use", "u":
			if len(blocks) > 1 {
				searchTerm := strings.Join(blocks[1:], " ")
				searchTerm = strings.Trim(searchTerm, "\"")
				result := preciseSearchModule(searchTerm)
				if len(result) == 1 {
					useID = result[0].ID
					currentPrefix = fmt.Sprintf("%s %s (%d_%s) > ", global.Name, result[0].Provider.Name, result[0].ID, result[0].Name)
				} else {
					logger.Println.Warn(fmt.Sprintf("没有找到 ID 为 %s 的模块", searchTerm))
				}
			} else {
				logger.Println.Warn("请指定要使用的模块。\n")
			}
		default:
			logger.Println.Warn("不存在该命令，请使用 help 命令查看帮助信息。\n")
		}
	}
}
