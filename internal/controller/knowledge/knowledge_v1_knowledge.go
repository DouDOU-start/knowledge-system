package knowledge

import (
	"context"
	v1 "knowledge-system-api/api/knowledge/v1"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

// ControllerV1 知识库V1接口控制器
// 只做参数校验、调用service、返回响应
// 错误全部用gerror/gcode统一处理
type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// BatchImport 批量导入知识条目
func (c *ControllerV1) BatchImport(ctx context.Context, req *v1.BatchImportReq) (res *v1.BatchImportRes, err error) {
	// 参数校验由框架自动完成，这里只需处理业务逻辑
	for _, item := range req.Items {
		// 调用LLM进行标签分类和摘要生成
		labels, summary, err := service.LLMClassifyByConfig(ctx, item.Content)
		if err != nil {
			g.Log().Errorf(ctx, "LLM推理失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "LLM推理失败: %s", err.Error())
		}

		// 过滤标签
		filtered := service.FilterLabels(labels, 3)

		// 向量化
		vector, err := service.Vectorize(ctx, item.Content)
		if err != nil {
			g.Log().Errorf(ctx, "向量化失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "向量化失败: %s", err.Error())
		}

		// 生成ID
		id := item.ID
		if id == "" {
			id = uuid.NewString()
		}

		// 存入向量数据库
		if err := service.QdrantUpsert(id, vector, item.Content, filtered, summary); err != nil {
			g.Log().Errorf(ctx, "Qdrant入库失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "Qdrant入库失败: %s", err.Error())
		}

		// 存入MySQL
		if err := service.KnowledgeService().CreateKnowledge(ctx, id, item.Content, filtered, summary); err != nil {
			g.Log().Errorf(ctx, "MySQL入库失败: %v", err)
			return nil, gerror.NewCodef(gcode.CodeInternalError, "MySQL入库失败: %s", err.Error())
		}
	}

	return &v1.BatchImportRes{Success: true}, nil
}

// Classify 单条内容标签打分
func (c *ControllerV1) Classify(ctx context.Context, req *v1.ClassifyReq) (res *v1.ClassifyRes, err error) {
	// 参数校验由框架自动完成

	// 调用LLM进行标签分类和摘要生成
	labels, summary, err := service.LLMClassifyByConfig(ctx, req.Content)
	if err != nil {
		g.Log().Errorf(ctx, "LLM推理失败: %v", err)
		return nil, gerror.NewCodef(gcode.CodeInternalError, "LLM推理失败: %s", err.Error())
	}

	// 转换为API响应格式
	var outLabels []v1.LabelScore
	for _, l := range labels {
		outLabels = append(outLabels, v1.LabelScore{
			LabelID: l.LabelID,
			Score:   l.Score,
		})
	}

	return &v1.ClassifyRes{
		Labels:  outLabels,
		Summary: summary,
	}, nil
}

// Search 知识检索
func (c *ControllerV1) Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error) {
	// 参数校验由框架自动完成

	// 设置默认值
	if req.TopK <= 0 {
		req.TopK = 5
	}

	// 根据模式选择不同的搜索方式
	var items []model.SearchResult
	switch req.Mode {
	case "keyword":
		items, err = service.KnowledgeService().SearchKnowledgeByKeyword(ctx, req.Query, req.TopK)
	case "semantic":
		items, err = service.KnowledgeService().SearchKnowledgeBySemantic(ctx, req.Query, req.TopK)
	case "hybrid", "":
		items, err = service.KnowledgeService().SearchKnowledgeByHybrid(ctx, req.Query, req.TopK)
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "不支持的搜索模式")
	}

	if err != nil {
		g.Log().Errorf(ctx, "知识检索失败: %v", err)
		return nil, gerror.NewCodef(gcode.CodeInternalError, "知识检索失败: %s", err.Error())
	}

	// 转换为API响应格式
	var outItems []v1.KnowledgeResult
	for _, item := range items {
		var outLabels []v1.LabelScore
		for _, l := range item.Labels {
			outLabels = append(outLabels, v1.LabelScore{
				LabelID: l.LabelID,
				Score:   l.Score,
			})
		}

		outItems = append(outItems, v1.KnowledgeResult{
			ID:      item.ID,
			Content: item.Content,
			Labels:  outLabels,
			Summary: item.Summary,
		})
	}

	return &v1.SearchRes{Items: outItems}, nil
}
