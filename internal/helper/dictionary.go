package helper

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

// =================================================================
// 1. 定义接口 (Interface)
// =================================================================

type IDictionary interface {
	GetID(ctx context.Context, label string) (id uint32, found bool)
}

// =================================================================
// 2. 服务注册与获取 (Service Registration & Retrieval) - 【这是之前缺失的部分】
// =================================================================

var (
	// localDictionary 是一个包内变量，用于存储 IDictionary 接口的单例实现。
	localDictionary IDictionary
)

// Dictionary 函数返回已注册的字典服务单例。
// 项目中其他模块通过调用 service.Dictionary() 来获取服务实例。
func Dictionary() IDictionary {
	if localDictionary == nil {
		// 这个 panic 会在忘记注册服务时提醒开发者
		panic("implement not found for interface IDictionary, forgot register?")
	}
	return localDictionary
}

// RegisterDictionary 函数用于将一个 IDictionary 的实现注册为单例服务。
func RegisterDictionary(i IDictionary) {
	localDictionary = i
}

// =================================================================
// 3. 服务实现 (Implementation)
// =================================================================

// sDictionary 是 IDictionary 接口的具体实现。
type sDictionary struct {
	labelToID map[string]uint32 // 内存中的只读映射表
}

// init 函数在包被导入时自动执行。
// 这是 Go 语言的特性，非常适合用来做服务注册。
func init() {
	// 在这里，我们将 New() 函数创建的实例注册为服务单例。
	// RegisterDictionary 函数现在已经被定义，所以这个调用不会再报错。
	RegisterDictionary(New())
}

// New 函数创建并返回一个具体的服务实现，它是 IDictionary 类型。
func New() IDictionary {
	// 从配置文件获取映射文件路径
	filePath := g.Cfg().MustGet(context.Background(), "dictionary.mapping_file").String()
	if filePath == "" {
		g.Log().Fatal(context.Background(), "字典映射文件路径 'dictionary.mapping_file' 未在配置文件中定义！")
	}

	// 加载并解析文件
	mapping, err := loadMappingFromFile(filePath)
	if err != nil {
		g.Log().Fatalf(context.Background(), "加载字典映射文件 '%s' 失败: %v", filePath, err)
	}

	g.Log().Infof(context.Background(), "成功从 '%s' 加载 %d 个标签映射", filePath, len(mapping))

	// 返回实现了 IDictionary 接口的 sDictionary 结构体实例
	return &sDictionary{
		labelToID: mapping,
	}
}

// GetID 是接口方法的具体实现
func (s *sDictionary) GetID(ctx context.Context, label string) (uint32, bool) {
	id, ok := s.labelToID[label]
	return id, ok
}

// loadMappingFromFile 是一个辅助函数，用于从JSON文件加载和合并映射
func loadMappingFromFile(path string) (map[string]uint32, error) {
	if !gfile.Exists(path) {
		return nil, gerror.Newf("映射文件不存在: %s", path)
	}

	content := gfile.GetContents(path)

	var structuredMap struct {
		C1_Topic map[string]uint32 `json:"C1_Topic"`
		C2_Type  map[string]uint32 `json:"C2_Type"`
	}

	if err := json.Unmarshal([]byte(content), &structuredMap); err != nil {
		return nil, gerror.Wrap(err, "解析JSON映射文件失败")
	}

	finalMap := make(map[string]uint32)
	for label, id := range structuredMap.C1_Topic {
		finalMap[label] = id
	}
	for label, id := range structuredMap.C2_Type {
		finalMap[label] = id
	}

	return finalMap, nil
}
