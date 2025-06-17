package knowledge

import (
	"context"
	"encoding/json"
	"knowledge-system-api/internal/dao"
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/model/do"
	"knowledge-system-api/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Knowledge 知识库业务逻辑实现
type Knowledge struct{}

// New 创建知识库业务逻辑实例
func New() *Knowledge {
	return &Knowledge{}
}

// CreateKnowledge 创建知识条目
func (s *Knowledge) CreateKnowledge(ctx context.Context, id, content string, labels []model.LabelScore, summary string) error {
	labelsJson, err := json.Marshal(labels)
	if err != nil {
		return err
	}

	now := gtime.Now()
	_, err = dao.Knowledge.Ctx(ctx).Data(do.Knowledge{
		Id:        id,
		Content:   content,
		Labels:    string(labelsJson),
		Summary:   summary,
		CreatedAt: now,
		UpdatedAt: now,
	}).InsertAndGetId()

	return err
}

// GetKnowledgeById 根据ID获取知识条目
func (s *Knowledge) GetKnowledgeById(ctx context.Context, id string) (*model.KnowledgeItem, error) {
	var entity entity.Knowledge
	err := dao.Knowledge.Ctx(ctx).Where(do.Knowledge{
		Id: id,
	}).Scan(&entity)
	if err != nil {
		return nil, err
	}

	if entity.Id == "" {
		return nil, nil
	}

	var labels []model.LabelScore
	if err := json.Unmarshal([]byte(entity.Labels), &labels); err != nil {
		g.Log().Warning(ctx, "解析标签JSON失败", err)
		labels = []model.LabelScore{}
	}

	return &model.KnowledgeItem{
		ID:        entity.Id,
		Content:   entity.Content,
		Labels:    labels,
		Summary:   entity.Summary,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

// SearchKnowledgeByKeyword 关键词搜索知识条目
func (s *Knowledge) SearchKnowledgeByKeyword(ctx context.Context, keyword string, limit int) ([]model.SearchResult, error) {
	var entities []entity.Knowledge

	// 使用 LIKE 进行关键词搜索
	err := dao.Knowledge.Ctx(ctx).
		WhereLike("content", "%"+keyword+"%").
		WhereOr("summary LIKE ?", "%"+keyword+"%").
		Limit(limit).
		OrderDesc("created_at").
		Scan(&entities)

	if err != nil {
		return nil, err
	}

	var results []model.SearchResult
	for _, e := range entities {
		var labels []model.LabelScore
		if err := json.Unmarshal([]byte(e.Labels), &labels); err != nil {
			g.Log().Warning(ctx, "解析标签JSON失败", err)
			labels = []model.LabelScore{}
		}

		results = append(results, model.SearchResult{
			ID:      e.Id,
			Content: e.Content,
			Labels:  labels,
			Summary: e.Summary,
			Score:   1.0, // 关键词搜索默认分数为1
		})
	}

	return results, nil
}

// SearchKnowledgeBySemantic 语义搜索知识条目
func (s *Knowledge) SearchKnowledgeBySemantic(ctx context.Context, query string, limit int) ([]model.SearchResult, error) {
	// 调用向量搜索服务
	if helper.Vectorize == nil {
		return nil, gerror.New("向量化服务未初始化")
	}

	vector, err := helper.Vectorize(ctx, query)
	if err != nil {
		return nil, err
	}

	// 调用 Qdrant 搜索
	if helper.VectorSearch == nil {
		return nil, gerror.New("向量搜索服务未初始化")
	}

	qdrantResults, err := helper.VectorSearch(query, vector, limit)
	if err != nil {
		return nil, err
	}

	var results []model.SearchResult
	for _, item := range qdrantResults {
		// 获取完整知识条目
		knowledgeItem, err := s.GetKnowledgeById(ctx, item.ID)
		if err != nil {
			g.Log().Warning(ctx, "获取知识条目失败", err)
			continue
		}

		if knowledgeItem == nil {
			continue
		}

		results = append(results, model.SearchResult{
			ID:      knowledgeItem.ID,
			Content: knowledgeItem.Content,
			Labels:  knowledgeItem.Labels,
			Summary: knowledgeItem.Summary,
			Score:   item.Score,
		})
	}

	return results, nil
}

// SearchKnowledgeByHybrid 混合搜索知识条目（关键词+语义）
func (s *Knowledge) SearchKnowledgeByHybrid(ctx context.Context, query string, limit int) ([]model.SearchResult, error) {
	// 先进行语义搜索
	semanticResults, err := s.SearchKnowledgeBySemantic(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// 再进行关键词搜索
	keywordResults, err := s.SearchKnowledgeByKeyword(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	// 合并结果，去重
	resultMap := make(map[string]model.SearchResult)

	// 先添加语义搜索结果
	for _, r := range semanticResults {
		resultMap[r.ID] = r
	}

	// 再添加关键词搜索结果，如果已存在则保留分数更高的
	for _, r := range keywordResults {
		if existing, ok := resultMap[r.ID]; ok {
			if r.Score > existing.Score {
				resultMap[r.ID] = r
			}
		} else {
			resultMap[r.ID] = r
		}
	}

	// 转换为数组
	var results []model.SearchResult
	for _, r := range resultMap {
		results = append(results, r)
	}

	// 限制返回数量
	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
