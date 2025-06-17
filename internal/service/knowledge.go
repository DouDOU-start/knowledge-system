package service

import (
	"knowledge-system-api/internal/service/interfaces"
)

// Knowledge 获取知识服务实例
func Knowledge() interfaces.KnowledgeService {
	return KnowledgeService()
}
