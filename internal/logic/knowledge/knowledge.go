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

// SearchKnowledgeByHybrid 混合搜索知识条目（基于用户意图的语义检索）
func (s *Knowledge) SearchKnowledgeByHybrid(ctx context.Context, query string, repoName string, limit uint64) ([]model.SearchResult, error) {
	g.Log().Debug(ctx, "开始混合搜索，基于标签和语义检索")

	// 步骤1：分析用户查询意图，提取关键标签
	var labelScores []model.LabelScore

	if helper.LLMClassify != nil {
		g.Log().Debug(ctx, "分析用户查询意图")
		var err error
		labelScores, _, err = helper.LLMClassify(ctx, query)
		if err != nil {
			g.Log().Warningf(ctx, "LLM分析失败: %v, 将使用纯向量搜索", err)
		}
	}

	// 使用高优先级标签进行过滤的向量检索
	g.Log().Debugf(ctx, "开始向量检索，标签数量: %d", len(labelScores))
	qdrantResults, err := helper.VectorSearch(repoName, query, labelScores, limit)
	if err != nil {
		return nil, err
	}

	// 步骤4：处理结果
	var results []model.SearchResult
	for _, item := range qdrantResults {
		// 获取完整知识条目
		knowledgeItem, err := s.GetKnowledgeById(ctx, item.ID)
		if err != nil || knowledgeItem == nil {
			g.Log().Warningf(ctx, "获取知识条目失败: %v", err)
			continue
		}

		// 如果指定了知识库名称，但结果不匹配，则跳过
		if repoName != "" && knowledgeItem.RepoName != repoName {
			continue
		}

		// 添加到结果集
		results = append(results, model.SearchResult{
			ID:       knowledgeItem.ID,
			RepoName: knowledgeItem.RepoName,
			Content:  knowledgeItem.Content,
			Labels:   knowledgeItem.Labels,
			Summary:  knowledgeItem.Summary,
			Score:    item.Score,
		})
	}

	// 按得分排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	g.Log().Debugf(ctx, "混合搜索完成: 共返回 %d 条结果", len(results))
	return results, nil
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
