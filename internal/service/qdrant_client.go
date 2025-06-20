package service

import (
	"context"
	"fmt"
	"knowledge-system-api/internal/helper"
	"knowledge-system-api/internal/model"
	"sort"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/qdrant/go-client/qdrant"
)

var (
	qdrantConfigInstance *QdrantConfig
	qdrantOnce           sync.Once
	qdrantClientOnce     sync.Once
	qdrantClient         *qdrant.Client
)

// QdrantConfig Qdrant配置
type QdrantConfig struct {
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	Dimension uint64 `yaml:"dimension" json:"dimension"`
}

// LoadQdrantConfig 加载Qdrant配置
func LoadQdrantConfig(ctx context.Context) (*QdrantConfig, error) {
	var cfg QdrantConfig
	if err := g.Cfg().MustGet(ctx, "qdrant").Scan(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetQdrantConfig 获取Qdrant配置单例
func GetQdrantConfig() *QdrantConfig {
	qdrantOnce.Do(func() {
		ctx := gctx.New()
		cfg, err := LoadQdrantConfig(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "加载qdrant配置失败: %v", err)
			// 使用默认配置
			cfg = &QdrantConfig{
				Host:      "localhost",
				Port:      6334,
				Dimension: 1024,
			}
		}
		qdrantConfigInstance = cfg
		g.Log().Infof(ctx, "Qdrant配置加载成功：%+v", qdrantConfigInstance)
	})
	return qdrantConfigInstance
}

// InitQdrantClient 初始化Qdrant客户端 - 公开方法，可在程序启动时调用
func InitQdrantClient(ctx context.Context) error {
	var initErr error
	qdrantClientOnce.Do(func() {
		// 确保配置已加载
		qdrantConfigInstance = GetQdrantConfig()

		config := &qdrant.Config{
			Host: qdrantConfigInstance.Host,
			Port: qdrantConfigInstance.Port,
		}

		g.Log().Infof(ctx, "正在连接Qdrant服务: %s:%d", config.Host, config.Port)

		// 创建客户端
		var err error
		qdrantClient, err = qdrant.NewClient(config)
		if err != nil {
			g.Log().Errorf(ctx, "Qdrant客户端初始化失败: %v", err)
			initErr = fmt.Errorf("qdrant客户端初始化失败: %w", err)
			return
		}

		// 测试连接
		collections, err := qdrantClient.ListCollections(ctx)
		if err != nil {
			g.Log().Errorf(ctx, "连接Qdrant服务失败: %v", err)
			qdrantClient = nil
			initErr = fmt.Errorf("连接Qdrant服务失败: %w", err)
			return
		}

		g.Log().Infof(ctx, "Qdrant客户端初始化成功，当前存在 %d 个集合", len(collections))
	})
	return initErr
}

// GetQdrantClient 获取Qdrant客户端实例
func GetQdrantClient(ctx context.Context) (*qdrant.Client, error) {
	if qdrantClient == nil {
		err := InitQdrantClient(ctx)
		if err != nil {
			return nil, err
		}
	}
	return qdrantClient, nil
}

// CloseQdrantClient 关闭Qdrant客户端连接
func CloseQdrantClient() {
	if qdrantClient != nil {
		g.Log().Info(gctx.New(), "正在关闭Qdrant客户端连接")
		qdrantClient = nil
	}
}

// DeleteQdrantCollection 删除指定的集合（仅用于开发和调试）
func DeleteQdrantCollection(ctx context.Context, collectionName string) error {
	client, err := GetQdrantClient(ctx)
	if err != nil {
		return fmt.Errorf("获取Qdrant客户端失败: %w", err)
	}

	exists, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("检查集合是否存在时出错: %w", err)
	}

	if !exists {
		g.Log().Infof(ctx, "集合 %s 不存在，无需删除", collectionName)
		return nil
	}

	err = client.DeleteCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("删除集合失败: %w", err)
	}

	g.Log().Infof(ctx, "集合 %s 已成功删除", collectionName)
	return nil
}

// QdrantUpsert 将知识条目写入Qdrant向量库
func QdrantUpsert(ctx context.Context, repoName string, id string, content string, summary string, labels []model.LabelScore) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 参数检查
	if repoName == "" {
		return fmt.Errorf("QdrantUpsert: 集合名称不能为空")
	}

	if id == "" {
		return fmt.Errorf("QdrantUpsert: ID不能为空")
	}

	// 获取客户端，如果不存在则初始化
	client, err := GetQdrantClient(ctx)
	if err != nil {
		return fmt.Errorf("QdrantUpsert: %w", err)
	}

	// 2. 准备标签数据
	var labelPoints []interface{}
	for _, l := range labels {
		labelPoints = append(labelPoints, map[string]interface{}{
			"label_id": l.Name,
			"score":    l.Score,
		})
	}

	// 3. 构建payload
	payload := qdrant.NewValueMap(map[string]any{
		"content": content,
		"summary": summary,
		"labels":  labelPoints,
	})

	// 4. 确保集合存在
	collectionExists, err := client.CollectionExists(ctx, repoName)
	if err != nil {
		g.Log().Warningf(ctx, "检查集合是否存在时出错: %v", err)
	}

	// 如果集合不存在，创建集合
	if !collectionExists {
		// 创建密集向量和稀疏向量配置
		vectorsConfig := qdrant.NewVectorsConfigMap(map[string]*qdrant.VectorParams{
			"content_dense": {
				Size:     qdrantConfigInstance.Dimension, // 密集向量维度
				Distance: qdrant.Distance_Cosine,
			},
		})
		sparseVectorsConfig := qdrant.NewSparseVectorsConfig(map[string]*qdrant.SparseVectorParams{
			"labels_sparse": {},
		})

		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName:      repoName,
			VectorsConfig:       vectorsConfig,
			SparseVectorsConfig: sparseVectorsConfig,
		})
		if err != nil {
			return fmt.Errorf("创建集合失败: %w", err)
		}
		g.Log().Infof(ctx, "成功创建集合 %s，包含向量: content_dense, labels_sparse", repoName)
	} else {
		g.Log().Debugf(ctx, "集合 %s 已存在，跳过创建", repoName)
	}

	// 生成密集向量
	denseVector, err := helper.Vectorize(ctx, content)

	if err != nil {
		return fmt.Errorf("向量化内容失败: %w", err)
	}

	// 生成稀疏向量
	var sparseIndices []uint32
	var sparseValues []float32

	for _, l := range labels {
		// 调用GetID方法
		id, found := helper.Dictionary().GetID(ctx, l.Name)
		if found {
			sparseIndices = append(sparseIndices, id)
			sparseValues = append(sparseValues, l.Score)
		} else {
			// 如果标签未在字典中找到，可以选择记录一个警告日志
			g.Log().Warningf(ctx, "标签 '%s' 在预定义字典中未找到，已忽略。", l.Name)
		}
	}

	vectorsMap := map[string]*qdrant.Vector{
		"content_dense": qdrant.NewVectorDense(denseVector),
		"labels_sparse": qdrant.NewVectorSparse(sparseIndices, sparseValues),
	}

	// 5. 上传点
	_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: repoName,
		Points: []*qdrant.PointStruct{
			{
				Id: qdrant.NewIDUUID(id),
				Vectors: &qdrant.Vectors{
					VectorsOptions: &qdrant.Vectors_Vectors{
						Vectors: &qdrant.NamedVectors{
							Vectors: vectorsMap,
						},
					},
				},
				Payload: payload,
			},
		},
		Wait: func() *bool { b := true; return &b }(), // 等待上传完成
	})

	if err != nil {
		return fmt.Errorf("上传向量到Qdrant失败: %w", err)
	}

	return nil
}

// QdrantSearch 向量搜索
func QdrantSearch(ctx context.Context, repoName string, content string, labels []model.LabelScore, limit uint64) ([]model.VectorSearchResult, error) {
	// 参数检查
	if repoName == "" {
		return nil, fmt.Errorf("QdrantSearch: 集合名称不能为空")
	}

	// 获取客户端，如果不存在则初始化
	client, err := GetQdrantClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("QdrantSearch: %w", err)
	}

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 生成稀疏向量 (用于预查询)
	// 使用map确保索引的唯一性
	sparseIndexMap := make(map[uint32]float32)

	for _, l := range labels {
		// 调用GetID方法
		id, found := helper.Dictionary().GetID(ctx, l.Name)

		// 如果在我们的预定义字典中找到了这个标签
		if found {
			// 如果该索引已存在，则取较大的分数
			if existingScore, exists := sparseIndexMap[id]; exists {
				if l.Score > existingScore {
					sparseIndexMap[id] = l.Score
				}
			} else {
				sparseIndexMap[id] = l.Score
			}
		} else {
			// 如果标签未在字典中找到，可以选择记录一个警告日志
			g.Log().Warningf(ctx, "标签 '%s' 在预定义字典中未找到，已忽略。", l.Name)
		}
	}

	// 将map转换为有序列表
	var sparseIndices []uint32
	var sparseValues []float32

	// 1. 收集所有键
	for idx := range sparseIndexMap {
		sparseIndices = append(sparseIndices, idx)
	}

	// 2. 排序键（Qdrant要求索引必须排序）
	sort.Slice(sparseIndices, func(i, j int) bool {
		return sparseIndices[i] < sparseIndices[j]
	})

	// 3. 按排序后的键收集值
	for _, idx := range sparseIndices {
		sparseValues = append(sparseValues, sparseIndexMap[idx])
	}

	// 日志记录稀疏向量信息
	if len(sparseIndices) > 0 {
		g.Log().Debugf(ctx, "搜索使用稀疏向量: 索引数量=%d", len(sparseIndices))
	} else {
		g.Log().Warningf(ctx, "未使用稀疏向量过滤，可能影响检索精度")
	}

	// 生成密集向量 (用于主查询)
	vector, err := helper.Vectorize(ctx, content)

	if err != nil {
		return nil, fmt.Errorf("向量化内容失败: %w", err)
	}

	// 定义第一阶段的 Prefetch 查询
	prefetchCoreQuery := &qdrant.Query{
		Variant: &qdrant.Query_Nearest{
			Nearest: &qdrant.VectorInput{
				Variant: &qdrant.VectorInput_Sparse{
					Sparse: &qdrant.SparseVector{
						Values:  sparseValues,
						Indices: sparseIndices,
					},
				},
			},
		},
	}

	prefetchLimit := limit * 3 // 预查询的限制，通常设置为主查询的3倍，以确保有足够的候选项
	prefetchQuery := &qdrant.PrefetchQuery{
		Query: prefetchCoreQuery, // 将上面定义的查询逻辑放入
		Limit: &prefetchLimit,    // 直接在 PrefetchQuery 层面设置 Limit
	}

	// 定义第二阶段的主查询 (基于稀疏向量)
	mainQuery := &qdrant.Query{
		Variant: &qdrant.Query_Nearest{
			Nearest: &qdrant.VectorInput{
				Variant: &qdrant.VectorInput_Dense{
					Dense: &qdrant.DenseVector{
						Data: vector,
					},
				},
			},
		},
	}

	queryPointsRequest := &qdrant.QueryPoints{
		CollectionName: repoName,
		Prefetch:       []*qdrant.PrefetchQuery{prefetchQuery}, // 放入预查询
		Query:          mainQuery,                              // 放入主查询
		Limit:          &limit,                                 // 最终限制
		WithPayload:    qdrant.NewWithPayload(true),            // 返回Payload
	}

	// 执行搜索
	results, err := client.Query(timeoutCtx, queryPointsRequest)
	if err != nil {
		g.Log().Errorf(ctx, "Qdrant搜索失败: %v", err)
		return nil, fmt.Errorf("qdrant搜索失败: %w", err)
	}

	// 处理结果
	var searchResults []model.VectorSearchResult
	for _, point := range results {
		// 获取payload
		payload := make(map[string]interface{})

		// 复制所有payload字段
		if point.Payload != nil {
			for k, v := range point.Payload {
				payload[k] = v
			}
		}

		searchResults = append(searchResults, model.VectorSearchResult{
			ID:      point.Id.String(),
			Score:   point.Score,
			Payload: payload,
		})
	}

	g.Log().Debugf(ctx, "Qdrant搜索完成，找到 %d 条结果", len(searchResults))
	return searchResults, nil
}
