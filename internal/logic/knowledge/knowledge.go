package knowledge

import (
	"context"
	"encoding/json"
	"knowledge-system-api/internal/dao"
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/model"
	"knowledge-system-api/internal/model/do"
	"knowledge-system-api/internal/model/entity"
	"sort"
	"strings"

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
func (s *Knowledge) CreateKnowledge(ctx context.Context, id, repoName, content string, labels []model.LabelScore, summary string) error {
	labelsJson, err := json.Marshal(labels)
	if err != nil {
		return err
	}

	now := gtime.Now()
	_, err = dao.Knowledge.Ctx(ctx).Data(do.Knowledge{
		Id:        id,
		RepoName:  repoName,
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
		RepoName:  entity.RepoName,
		Content:   entity.Content,
		Labels:    labels,
		Summary:   entity.Summary,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

// SearchKnowledgeByKeyword 关键词搜索知识条目
func (s *Knowledge) SearchKnowledgeByKeyword(ctx context.Context, keyword string, repoName string, limit int) ([]model.SearchResult, error) {
	var entities []entity.Knowledge

	// 构建查询
	query := dao.Knowledge.Ctx(ctx).
		WhereLike("content", "%"+keyword+"%").
		WhereOr("summary LIKE ?", "%"+keyword+"%")

	// 如果指定了知识库名称，则添加条件
	if repoName != "" {
		query = query.Where("repo_name", repoName)
	}

	// 执行查询
	err := query.Limit(limit).
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
			ID:       e.Id,
			RepoName: e.RepoName,
			Content:  e.Content,
			Labels:   labels,
			Summary:  e.Summary,
			Score:    1.0, // 关键词搜索默认分数为1
		})
	}

	return results, nil
}

// SearchKnowledgeBySemantic 语义搜索知识条目
func (s *Knowledge) SearchKnowledgeBySemantic(ctx context.Context, query string, repoName string, limit int) ([]model.SearchResult, error) {
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

	// 这里传入 repoName 参数
	qdrantResults, err := helper.VectorSearch(query, vector, repoName, limit)
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

		// 如果指定了知识库名称，但结果不匹配，则跳过
		if repoName != "" && knowledgeItem.RepoName != repoName {
			continue
		}

		results = append(results, model.SearchResult{
			ID:       knowledgeItem.ID,
			RepoName: knowledgeItem.RepoName,
			Content:  knowledgeItem.Content,
			Labels:   knowledgeItem.Labels,
			Summary:  knowledgeItem.Summary,
			Score:    item.Score,
		})
	}

	return results, nil
}

// SearchKnowledgeByHybrid 混合搜索知识条目（基于用户意图的语义检索）
func (s *Knowledge) SearchKnowledgeByHybrid(ctx context.Context, query string, repoName string, limit int) ([]model.SearchResult, error) {
	g.Log().Debug(ctx, "开始混合搜索，基于用户意图的语义检索")

	// 步骤1：先分析用户查询意图，提取关键标签
	// 使用LLM对用户查询进行意图分析，提取出关键标签
	var targetLabels []string
	var targetLabelScores map[string]int = make(map[string]int)
	var enhancedQuery string = query // 默认使用原始查询

	if helper.LLMClassify != nil {
		g.Log().Debug(ctx, "使用LLM分析用户查询意图")
		labelScores, summary, err := helper.LLMClassify(ctx, query)
		if err == nil && len(labelScores) > 0 {
			// 只保留分数超过3分的标签
			var filteredLabels []model.LabelScore
			for _, ls := range labelScores {
				if ls.Score > 3 { // 过滤低分标签
					filteredLabels = append(filteredLabels, ls)
					targetLabels = append(targetLabels, ls.LabelID)
					targetLabelScores[ls.LabelID] = ls.Score
				}
			}
			if len(filteredLabels) > 0 {
				g.Log().Debugf(ctx, "从查询中识别出的有效标签(分数>3): %v", targetLabels)

				// 增强查询：将高分标签和原始查询组合
				highPriorityLabels := getHighPriorityLabels(targetLabels, targetLabelScores, 2)
				if len(highPriorityLabels) > 0 {
					// 将高优先级标签添加到查询中，以增强语义搜索效果
					enhancedQuery = query + " " + strings.Join(highPriorityLabels, " ")
					g.Log().Debugf(ctx, "增强后的查询: %s", enhancedQuery)
				}
			} else {
				g.Log().Debug(ctx, "未识别出分数超过3分的标签")
			}

			// 如果有摘要，使用摘要进一步理解查询意图
			if summary != "" {
				g.Log().Debugf(ctx, "查询意图摘要: %s", summary)
				// 可以使用摘要进一步优化查询
			}
		} else if err != nil {
			g.Log().Warningf(ctx, "LLM分类失败: %v，将使用原始查询", err)
		}
	} else {
		g.Log().Warning(ctx, "LLM分类服务未初始化，将使用原始查询")
	}

	// 步骤2：向量化查询（使用增强后的查询）
	if helper.Vectorize == nil {
		return nil, gerror.New("向量化服务未初始化")
	}

	// 对用户查询进行向量化（使用增强后的查询）
	vector, err := helper.Vectorize(ctx, enhancedQuery)
	if err != nil {
		return nil, err
	}

	// 步骤3：调用 Qdrant 搜索（使用增强后的查询和向量）
	if helper.VectorSearch == nil {
		return nil, gerror.New("向量搜索服务未初始化")
	}

	// 使用增强后的查询进行向量搜索
	// 如果有高分标签，可以启用基于标签的过滤
	g.Log().Debugf(ctx, "使用增强查询进行语义检索")
	qdrantResults, err := helper.VectorSearch(enhancedQuery, vector, repoName, limit)
	if err != nil {
		return nil, err
	}

	// 步骤4：处理搜索结果
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

		// 如果指定了知识库名称，但结果不匹配，则跳过
		if repoName != "" && knowledgeItem.RepoName != repoName {
			continue
		}

		// 创建结果项
		result := model.SearchResult{
			ID:       knowledgeItem.ID,
			RepoName: knowledgeItem.RepoName,
			Content:  knowledgeItem.Content,
			Labels:   knowledgeItem.Labels,
			Summary:  knowledgeItem.Summary,
			Score:    item.Score, // 使用向量搜索得分
		}

		// 添加一些调试信息
		g.Log().Debugf(ctx, "匹配项: ID=%s, 分数=%.4f", knowledgeItem.ID, item.Score)

		results = append(results, result)
	}

	// 步骤5：关键词补充（如果结果不足）
	if len(results) < limit {
		g.Log().Debugf(ctx, "语义搜索结果不足，使用关键词搜索补充")
		keywordResults, err := s.SearchKnowledgeByKeyword(ctx, query, repoName, limit-len(results))
		if err == nil && len(keywordResults) > 0 {
			// 创建已有结果的ID映射，避免重复
			existingIds := make(map[string]bool)
			for _, r := range results {
				existingIds[r.ID] = true
			}

			// 添加非重复的关键词搜索结果
			for _, kr := range keywordResults {
				if _, exists := existingIds[kr.ID]; !exists {
					// 关键词搜索的结果降低一些权重
					kr.Score = kr.Score * 0.8
					results = append(results, kr)
				}
			}
		}
	}

	// 按得分排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score // 降序排列
	})

	// 限制返回数量
	if len(results) > limit {
		results = results[:limit]
	}

	g.Log().Debugf(ctx, "混合搜索完成: 共返回 %d 条结果", len(results))

	return results, nil
}

// getHighPriorityLabels 获取高优先级标签（分数最高的前N个标签）
func getHighPriorityLabels(labels []string, scores map[string]int, topN int) []string {
	if len(labels) <= topN {
		return labels // 如果标签总数小于等于topN，直接返回全部
	}

	// 创建标签-分数对
	type labelScore struct {
		label string
		score int
	}

	var pairs []labelScore
	for _, label := range labels {
		pairs = append(pairs, labelScore{label: label, score: scores[label]})
	}

	// 按分数降序排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	// 提取前topN个标签
	result := make([]string, 0, topN)
	for i := 0; i < topN && i < len(pairs); i++ {
		result = append(result, pairs[i].label)
	}

	return result
}

// GetAllRepos 获取所有知识库名称
func (s *Knowledge) GetAllRepos(ctx context.Context) ([]string, error) {
	var repos []string

	// 查询所有不同的知识库名称
	err := dao.Knowledge.Ctx(ctx).
		Fields("DISTINCT repo_name").
		OrderAsc("repo_name").
		Scan(&repos)

	if err != nil {
		return nil, err
	}

	return repos, nil
}

// RecoverTasks 恢复未完成的任务
// 在服务启动时调用
func (s *Knowledge) RecoverTasks() {
	ctx := context.Background()
	g.Log().Info(ctx, "开始恢复未完成任务...")

	// 查询所有未完成的任务
	var tasks []entity.ImportTask
	err := dao.ImportTask.Ctx(ctx).
		WhereIn("status", []string{"pending", "processing"}).
		Scan(&tasks)

	if err != nil {
		g.Log().Error(ctx, "查询未完成任务失败:", err)
		return
	}

	if len(tasks) == 0 {
		g.Log().Info(ctx, "没有发现未完成任务")
		return
	}

	g.Log().Infof(ctx, "发现 %d 个未完成任务，开始恢复...", len(tasks))

	// 恢复每个任务
	for _, task := range tasks {
		taskId := task.Id
		g.Log().Infof(ctx, "恢复任务 %s, 当前状态: %s, 进度: %d%%",
			taskId, task.Status, task.Progress)

		// 这里可以根据任务类型调用不同的处理逻辑
		// 例如，可以根据保存在数据库中的任务项重新启动处理
		// ...

		// 标记任务为已恢复状态
		_, err := dao.ImportTask.Ctx(ctx).
			Data(do.ImportTask{
				Status:  "processing",
				Message: "服务重启后恢复",
			}).
			Where(do.ImportTask{Id: taskId}).
			Update()

		if err != nil {
			g.Log().Errorf(ctx, "更新任务 %s 状态失败: %v", taskId, err)
		} else {
			g.Log().Infof(ctx, "任务 %s 已标记为已恢复状态", taskId)
		}
	}

	g.Log().Info(ctx, "未完成任务恢复处理完毕")
}
