package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/wgpsec/cloudsword/utils"
	"github.com/wgpsec/cloudsword/utils/global"
	"github.com/wgpsec/cloudsword/utils/logger"
	"io"
	"os"
	"strings"
)

func capturePrint(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f() // 调用包含 fmt.Print 的函数
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func fnModel(m global.Module) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		for _, b := range global.GetBasicOptionsWithId(m.ID) {
			v, ok := req.Params.Arguments[b.Key].(string)
			if !ok && b.Required {
				return nil, fmt.Errorf("invalid name parameter%s", b.Key)
			}
			if b.Key == global.Detail {
				value := strings.ToLower(v)
				if utils.Contains([]string{"0", "false", "f", "no", "n"}, value) {
					v = global.False
				} else if utils.Contains([]string{"1", "true", "t", "yes", "y"}, value) {
					v = global.True
				} else {
					logger.Println.Error(fmt.Sprintf("%s 的值类型错误，请设置成 True 或者 False。", b.Key))
				}
			}
			global.UpdateBasicOptionValue(b.Key, v)
			fmt.Println(b.Key + " ==> " + v)
		}
		output := capturePrint(func() {
			runModule(getModuleByID(m.ID))
		})

		return mcp.NewToolResultText(output), nil
	}
}

func MCPServer(transport string, addr string) {
	s := server.NewMCPServer(
		global.Name,
		global.Version,
	)
	modes := Modules()
	for _, mode := range modes {
		opts := make([]mcp.ToolOption, 0)
		// 添加调用model的描述
		opts = append(opts, mcp.WithDescription(mode.Desc))
		// 批量增加调用model需要的参数的名称、描述、是否必须
		for _, opt := range mode.BasicOptions {
			optTemp := []mcp.PropertyOption{mcp.Description(opt.Introduce)}
			if opt.Required {
				optTemp = append(optTemp, mcp.Required())
			}
			opts = append(opts, mcp.WithString(opt.Key, optTemp...))
		}
		// 批量把调用的参数和模型绑定函数添加
		s.AddTool(mcp.NewTool(mode.Introduce, opts...), fnModel(mode))
	}

	if transport == "sse" {
		port := strings.Split(addr, ":")[2]
		sseServer := server.NewSSEServer(s, server.WithBaseURL(addr))
		logger.Println.Infof("SSE server listening on :%s", addr)
		if err := sseServer.Start(":" + port); err != nil {
			logger.Println.Fatalf("Server error: %v\n", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			logger.Println.Fatalf("Server error: %v\n", err)
		}
	}

}
