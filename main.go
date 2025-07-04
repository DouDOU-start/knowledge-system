package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"knowledge-system-api/internal/cmd"
	"knowledge-system-api/internal/helper"
	_ "knowledge-system-api/internal/logic"
	_ "knowledge-system-api/internal/packed"
)

// @title       知识库检索系统API
// @version     1.0
// @description 基于GoFrame的知识库检索系统API服务
// @schemes     http https
// @contact.url https://github.com/DouDOU-start/knowledge-system
func main() {
	// 执行所有初始化
	helper.InitAll()

	// 运行主命令
	cmd.Main.Run(gctx.GetInitCtx())
}
