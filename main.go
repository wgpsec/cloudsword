package main

import (
	"github.com/wgpsec/cloudsword/cmd"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "sse" {
			cmd.MCPServer()
		}
	} else {
		cmd.Run()
	}

}
