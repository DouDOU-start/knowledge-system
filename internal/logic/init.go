// Package logic 业务逻辑层
package logic

import (
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/logic/feedback"
	"knowledge-system-api/internal/logic/knowledge"
	"knowledge-system-api/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// InitServiceLogic 初始化服务层的业务逻辑实现
func InitServiceLogic() {
	// 初始化知识库服务的业务逻辑
	k := knowledge.New()
	service.RegisterKnowledgeLogic(
		k.CreateKnowledge,
		k.GetKnowledgeById,
		// k.SearchKnowledgeByKeyword,
		// k.SearchKnowledgeBySemantic,
		k.SearchKnowledgeByHybrid,
		k.CreateImportTask,
		k.GetTaskStatus,
		k.UpdateTaskStatus,
		k.GetAllRepos,
		k.RecoverTasks,
	)

	// 初始化反馈服务的业务逻辑
	f := feedback.New()
	service.RegisterFeedback(f)
}

func init() {
	// 注册初始化和清理函数
	helper.SetInitFunctions(helper.InitFunctions{
		InitServices:    initServices,
		CleanupServices: cleanupServices,
	})
}

// initServices 初始化所有服务
func initServices() {
	ctx := gctx.New()
	g.Log().Info(ctx, "开始初始化服务...")

	// 初始化业务逻辑实现
	InitServiceLogic()

	// 初始化Qdrant客户端
	if err := service.InitQdrantClient(ctx); err != nil {
		g.Log().Errorf(ctx, "Qdrant客户端初始化失败: %v", err)
	} else {
		g.Log().Info(ctx, "Qdrant客户端初始化成功")
	}

	// 在这里添加其他服务的初始化

	g.Log().Info(ctx, "服务初始化完成")
}

// cleanupServices 清理所有服务资源
func cleanupServices() {
	ctx := gctx.New()
	g.Log().Info(ctx, "开始清理服务资源...")

	// 关闭Qdrant客户端连接
	service.CloseQdrantClient()

	// 在这里添加其他服务的资源清理

	g.Log().Info(ctx, "服务资源清理完成")
}
