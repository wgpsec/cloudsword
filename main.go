package main

import (
	"github.com/wgpsec/cloudsword/cmd"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		addr := "http://localhost:8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
		}
		cmd.MCPServer(os.Args[1], addr)
	} else {
		cmd.Run()
	}

}
