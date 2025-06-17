package helper

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// GetLabelThreshold 获取标签过滤阈值
// 从配置文件读取阈值，如果未配置则使用默认值70
func GetLabelThreshold(ctx context.Context) int {
	threshold := g.Cfg().MustGet(ctx, "llm.label_threshold", 3).Int()
	return threshold
}
