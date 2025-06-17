# 知识库检索系统

基于 GoFrame 框架开发的知识库检索系统 API 服务，支持知识条目的批量导入、标签分类和多模式检索。

## 功能特性

- 知识条目批量导入
- 内容自动标签分类和摘要生成
- 多模式检索：关键词、语义、混合
- 向量数据库集成
- RESTful API 接口

## 技术栈

- GoFrame 2.x
- MySQL 数据库
- Qdrant 向量数据库
- Ollama/LangChain 支持

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+
- Qdrant 向量数据库
- Ollama (可选，用于本地嵌入和 LLM 推理)

### 安装与运行

1. 克隆代码库

```bash
git clone https://github.com/DouDOU-start/knowledge-system.git
cd knowledge-system
```

2. 修改配置文件

编辑 `config.yaml` 文件，配置数据库连接、向量数据库和模型参数。

3. 初始化数据库

```bash
# 执行 SQL 脚本创建数据库和表
mysql -u root -p < sql/init.sql
```

4. 编译运行

```bash
go mod tidy
go build -o server
./server
```

5. 访问 API 文档

浏览器访问 `http://localhost:8080/swagger` 查看 API 文档。

## API 接口

- `POST /api/v1/knowledge/batch_import` - 批量导入知识条目
- `POST /api/v1/knowledge/classify` - 单条内容标签打分
- `POST /api/v1/knowledge/search` - 知识检索

## 目录结构

```
.
├── api                 # API 接口定义
│   └── knowledge       # 知识库模块接口
├── internal            # 内部实现
│   ├── cmd             # 命令行入口
│   ├── controller      # 控制器层
│   ├── dao             # 数据访问层
│   ├── logic           # 业务逻辑层
│   ├── model           # 数据模型
│   └── service         # 服务接口层
├── resource            # 资源文件
│   └── public          # 公共资源
├── sql                 # SQL 脚本
├── config.yaml         # 配置文件
└── main.go             # 主程序入口
```

## 许可证

MIT