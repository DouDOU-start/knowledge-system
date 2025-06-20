package service

// import (
// 	"context"
// 	"knowledge-system-api/internal/helper"

// 	"github.com/gogf/gf/v2/frame/g"
// )

// func init() {
// 	helper.SetInitFunctions(helper.InitFunctions{
// 		InitServices: func() {
// 			ctx := context.Background()
// 			g.Log().Info(ctx, "初始化服务...")

// 			// 确保Qdrant客户端正确初始化
// 			cfg := GetQdrantConfig()
// 			g.Log().Infof(ctx, "Qdrant服务配置: %s:%d", cfg.Host, cfg.Port)

// 			// 准备恢复未完成任务
// 			RecoverUnfinishedTasks(ctx)
// 		},
// 		CleanupServices: func() {
// 			ctx := context.Background()
// 			g.Log().Info(ctx, "清理服务资源...")

// 			// 关闭Qdrant客户端连接
// 			CloseQdrantClient()

// 			g.Log().Info(ctx, "服务资源清理完成")
// 		},
// 	})
// }
