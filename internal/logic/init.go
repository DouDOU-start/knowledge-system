// Package logic 业务逻辑层
package logic

import (
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/logic/feedback"
	"knowledge-system-api/internal/logic/knowledge"
	"knowledge-system-api/internal/service"
)

// InitServiceLogic 初始化服务层的业务逻辑实现
func InitServiceLogic() {
	// 初始化知识库服务的业务逻辑
	k := knowledge.New()
	service.RegisterKnowledgeLogic(
		k.CreateKnowledge,
		k.GetKnowledgeById,
		k.SearchKnowledgeByKeyword,
		k.SearchKnowledgeBySemantic,
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
	// 注册初始化函数
	helper.SetInitFunctions(helper.InitFunctions{
		InitServices: InitServiceLogic,
	})
}
