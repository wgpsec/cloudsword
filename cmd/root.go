package cmd

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"github.com/wgpsec/cloudsword/utils"
)

func Run() {
	color.Green(banner)
	fmt.Println(projectInfo)
	utils.CheckVersion()

	p = prompt.New(
		executor,
		completer,
		prompt.OptionLivePrefix(changingPrefixLivePrefix),
	)
	p.Run()
}
