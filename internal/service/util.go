package service

import (
	"fmt"
	"regexp"
	// "regexp"
)

// CalcLabelScore 标签分数加权和
// func CalcLabelScore(queryLabels, itemLabels []model.LabelScore) float32 {
// 	var score float32
// 	for _, ql := range queryLabels {
// 		for _, il := range itemLabels {
// 			if ql.LabelID == il.LabelID {
// 				score += float32(ql.Score * il.Score)
// 			}
// 		}
// 	}
// 	return score
// }

// ExtractJSONFromLLMResponse 提取LLM返回中的JSON片段
func ExtractJSONFromLLMResponse(resp string) (string, error) {
	re := regexp.MustCompile("(?s)```json\\s*(\\{.*?\\})\\s*```")
	matches := re.FindStringSubmatch(resp)
	if len(matches) > 1 {
		return matches[1], nil
	}
	// 退而求其次，找第一个 { 到最后一个 }
	start, end := -1, -1
	for i, c := range resp {
		if c == '{' && start == -1 {
			start = i
		}
		if c == '}' {
			end = i
		}
	}
	if start != -1 && end != -1 && end > start {
		return resp[start : end+1], nil
	}
	return "", fmt.Errorf("未找到有效的JSON内容")
}
