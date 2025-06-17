package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"

	"knowledge-system-api/internal/cmd"
	_ "knowledge-system-api/internal/logic"
	_ "knowledge-system-api/internal/packed"
)

// @title       知识库检索系统API
// @version     1.0
// @description 基于GoFrame的知识库检索系统API服务
// @schemes     http https
// @contact.url https://github.com/your-username/knowledge-system
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
