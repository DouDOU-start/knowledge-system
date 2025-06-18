package feedback

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"knowledge-system-api/internal/dao"
	"knowledge-system-api/internal/model/do"
	"knowledge-system-api/internal/service"
)

// 标签结构体
type Label struct {
	Tag   string  `json:"tag"`
	Score float64 `json:"score"`
}

// 默认测试会话ID
const DefaultTestSessionID = "90c91010-d24b-4056-9c46-89aafe0ed4cb"

// Feedback 是反馈服务实现
type Feedback struct{}

// New 创建反馈服务
func New() *Feedback {
	return &Feedback{}
}

// Init 初始化
func (s *Feedback) Init(ctx context.Context) {
	service.RegisterFeedback(s)

	// 注册定时任务，按照GoFrame规范获取配置
	// cronPattern := g.Cfg().MustGet(ctx, "feedback.update.cron", "0 3 * * *").String()
	// _, err := g.Cron().Add(ctx, cronPattern, func(ctx context.Context) {
	// 	s.ProcessFeedbacks(ctx)
	// })

	// if err != nil {
	// 	g.Log().Error(ctx, "注册反馈处理定时任务失败:", err)
	// }
}

// Add 添加反馈
func (s *Feedback) Add(ctx context.Context, sessionID, userQuery, knowledgeID, action string) (int64, error) {
	// 如果没有提供会话ID，使用默认测试ID
	if sessionID == "" {
		sessionID = DefaultTestSessionID
	}

	// 检查知识ID是否存在
	knowledge, err := g.DB().Model("knowledge").Ctx(ctx).Where("id", knowledgeID).One()
	if err != nil {
		return 0, err
	}
	if knowledge.IsEmpty() {
		return 0, gerror.New("知识ID不存在")
	}

	// 保存反馈
	id, err := dao.Feedback.Ctx(ctx).Data(do.Feedback{
		SessionId:            sessionID,
		UserQuery:            userQuery,
		RetrievedKnowledgeId: knowledgeID,
		Action:               action,
		Timestamp:            gtime.Now(),
	}).InsertAndGetId()

	return id, err
}

// List 查询反馈列表
func (s *Feedback) List(ctx context.Context, sessionID, knowledgeID, action, startTime, endTime string, page, pageSize int) ([]service.FeedbackItem, int, error) {
	// 构建查询条件
	model := g.DB().Model("feedback").Ctx(ctx).Safe()

	if sessionID != "" {
		model = model.Where("session_id", sessionID)
	}

	if knowledgeID != "" {
		model = model.Where("retrieved_knowledge_id", knowledgeID)
	}

	if action != "" {
		model = model.Where("action", action)
	}

	if startTime != "" {
		model = model.Where("timestamp >= ?", startTime)
	}

	if endTime != "" {
		model = model.Where("timestamp <= ?", endTime)
	}

	// 查询总数
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	var feedbacks []map[string]interface{}
	err = model.Page(page, pageSize).Order("timestamp DESC").Scan(&feedbacks)
	if err != nil {
		return nil, 0, err
	}

	// 转换结果
	result := make([]service.FeedbackItem, len(feedbacks))
	for i, item := range feedbacks {
		result[i] = service.FeedbackItem{
			ID:          gconv.String(item["id"]),
			SessionID:   gconv.String(item["session_id"]),
			UserQuery:   gconv.String(item["user_query"]),
			KnowledgeID: gconv.String(item["retrieved_knowledge_id"]),
			Action:      gconv.String(item["action"]),
			Timestamp:   gconv.Time(item["timestamp"]),
		}
	}

	return result, total, nil
}

// ProcessFeedbacks 处理反馈数据，更新标签分数
func (s *Feedback) ProcessFeedbacks(ctx context.Context) error {
	g.Log().Info(ctx, "开始处理用户反馈，更新知识标签分数")

	// 按照GoFrame规范获取配置
	likeAdjust := g.Cfg().MustGet(ctx, "feedback.tag_adjustment.like", 0.05).Float64()
	dislikeAdjust := g.Cfg().MustGet(ctx, "feedback.tag_adjustment.dislike", -0.08).Float64()
	minScore := g.Cfg().MustGet(ctx, "feedback.tag_adjustment.min_score", 0.1).Float64()
	maxScore := g.Cfg().MustGet(ctx, "feedback.tag_adjustment.max_score", 10.0).Float64()

	// 查询过去24小时的反馈数据
	yesterday := gtime.Now().AddDate(0, 0, -1)
	var feedbacks []struct {
		KnowledgeID string `json:"retrieved_knowledge_id"`
		Action      string `json:"action"`
	}

	err := g.DB().Model("feedback").Ctx(ctx).
		Fields("retrieved_knowledge_id, action").
		Where("timestamp >= ?", yesterday).
		Scan(&feedbacks)
	if err != nil {
		return err
	}

	// 按知识ID分组统计反馈
	feedbackStats := make(map[string]map[string]int)
	for _, f := range feedbacks {
		if _, exists := feedbackStats[f.KnowledgeID]; !exists {
			feedbackStats[f.KnowledgeID] = map[string]int{
				"like":    0,
				"dislike": 0,
			}
		}
		feedbackStats[f.KnowledgeID][f.Action]++
	}

	// 处理每个知识的反馈
	for knowledgeID, stats := range feedbackStats {
		// 只有存在"dislike"或大量"like"才进行处理
		if stats["dislike"] == 0 && stats["like"] < 5 {
			continue
		}

		// 查询知识
		var knowledge struct {
			ID     string `json:"id"`
			Labels string `json:"labels"`
		}

		err := g.DB().Model("knowledge").Ctx(ctx).
			Fields("id, labels").
			Where("id", knowledgeID).
			Scan(&knowledge)
		if err != nil || knowledge.ID == "" {
			g.Log().Warningf(ctx, "获取知识失败: %v, ID: %s", err, knowledgeID)
			continue
		}

		// 解析标签
		var labels []Label
		if err := json.Unmarshal([]byte(knowledge.Labels), &labels); err != nil {
			g.Log().Warningf(ctx, "解析标签失败: %v, ID: %s", err, knowledgeID)
			continue
		}

		if len(labels) == 0 {
			continue
		}

		// 计算调整量
		adjustment := float64(stats["like"])*likeAdjust + float64(stats["dislike"])*dislikeAdjust

		// 记录是否有变化
		changed := false

		// 调整所有标签的分数
		for i := range labels {
			oldScore := labels[i].Score

			// 应用调整
			newScore := oldScore + adjustment

			// 确保分数在合理范围内
			if newScore < minScore {
				newScore = minScore
			} else if newScore > maxScore {
				newScore = maxScore
			}

			// 只有分数有变化才更新
			if newScore != oldScore {
				labels[i].Score = newScore
				changed = true
			}
		}

		// 有变化才更新数据库
		if changed {
			// 序列化标签
			labelsJson, err := json.Marshal(labels)
			if err != nil {
				g.Log().Warningf(ctx, "序列化标签失败: %v, ID: %s", err, knowledgeID)
				continue
			}

			// 更新知识的labels字段
			_, err = g.DB().Model("knowledge").Ctx(ctx).
				Data(g.Map{"labels": string(labelsJson)}).
				Where("id", knowledgeID).
				Update()

			if err != nil {
				g.Log().Errorf(ctx, "更新知识标签失败: %v, ID: %s", err, knowledgeID)
			} else {
				g.Log().Infof(ctx, "已更新知识标签, ID: %s, 点赞: %d, 点踩: %d, 调整值: %.2f",
					knowledgeID, stats["like"], stats["dislike"], adjustment)
			}
		}
	}

	g.Log().Info(ctx, "反馈处理完成")
	return nil
}
