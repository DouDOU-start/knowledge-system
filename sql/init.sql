-- =================================================================
-- 知识库系统数据库完整脚本 (最终优化版)
-- 包含: knowledge, import_task, import_task_item, task_queue 四张表
-- 核心优化:
-- 1. `import_task` 表中的 `items` 字段被拆分为独立的 `import_task_item` 表，实现结构规范化。
-- 2. 所有表结构一次性定义，避免后期 ALTER TABLE 操作。
-- 3. 状态字段 (`status`) 均使用 ENUM 类型进行优化。
-- 4. 索引经过设计，避免冗余并提高核心查询场景的性能。
-- 5. 自动更新时间戳字段 (`updated_at`) 均已配置。
-- =================================================================

-- 步骤 1: 创建数据库 (如果不存在)
CREATE DATABASE IF NOT EXISTS `knowledge_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 步骤 2: 使用数据库
USE `knowledge_system`;

-- 步骤 3: 创建知识表 (优化版)
CREATE TABLE IF NOT EXISTS `knowledge` (
  `id` varchar(36) NOT NULL COMMENT '唯一ID，服务端生成UUID',
  `repo_name` varchar(100) NOT NULL DEFAULT 'default' COMMENT '知识库名称',
  `content` text NOT NULL COMMENT '知识内容',
  `labels` json DEFAULT NULL COMMENT '标签分数数组',
  `summary` text DEFAULT NULL COMMENT '内容摘要',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_repo_name` (`repo_name`),
  FULLTEXT KEY `idx_content` (`content`) COMMENT '内容全文索引',
  FULLTEXT KEY `idx_summary` (`summary`) COMMENT '摘要全文索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识条目表';


-- 步骤 4: 创建导入任务表 (优化版 - 移除 items 字段)
CREATE TABLE IF NOT EXISTS `import_task` (
  `id` varchar(36) NOT NULL COMMENT '任务ID',
  `status` ENUM('pending', 'processing', 'completed', 'failed', 'completed_with_errors') NOT NULL DEFAULT 'pending' COMMENT '任务状态',
  `progress` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '处理进度，0-100',
  `total` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '总条目数',
  `processed` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '已处理条目数',
  `failed` int UNSIGNED NOT NULL DEFAULT 0 COMMENT '失败条目数',
  `message` varchar(255) DEFAULT NULL COMMENT '任务相关信息',
  -- `items` 字段已被移除
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status_created_at` (`status`, `created_at`) COMMENT '状态与创建时间复合索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导入任务表';


-- 步骤 5: 创建导入任务条目表 (新增)
CREATE TABLE IF NOT EXISTS `import_task_item` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '条目自增ID',
  `task_id` varchar(36) NOT NULL COMMENT '所属任务ID',
  `status` ENUM('pending', 'processing', 'completed', 'failed') NOT NULL DEFAULT 'pending' COMMENT '条目处理状态',
  `source_data` json NOT NULL COMMENT '原始数据 (如单条知识的JSON)',
  `error_message` text DEFAULT NULL COMMENT '处理失败时的错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id_status` (`task_id`, `status`) COMMENT '任务ID与状态复合索引，用于快速统计和查询特定任务下的条目'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='导入任务的单个条目表';


-- 步骤 6: 创建任务队列表 (已优化)
CREATE TABLE IF NOT EXISTS `task_queue` (
  `id` varchar(36) NOT NULL COMMENT '队列项ID',
  `task_id` varchar(36) NOT NULL COMMENT '任务ID',
  `priority` int NOT NULL DEFAULT 0 COMMENT '优先级',
  `status` ENUM('waiting', 'processing', 'completed', 'failed') NOT NULL DEFAULT 'waiting' COMMENT '状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `started_at` datetime DEFAULT NULL COMMENT '开始处理时间',
  `ended_at` datetime DEFAULT NULL COMMENT '处理结束时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_waiting_tasks` (`status`, `priority`, `created_at`) COMMENT '用于快速查询待处理任务'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务队列表';


-- 步骤 7: 创建用户并授权 (可选，根据实际情况修改)
-- CREATE USER 'knowledge_user'@'%' IDENTIFIED BY 'knowledge_password';
-- GRANT ALL PRIVILEGES ON knowledge_system.* TO 'knowledge_user'@'%';
-- FLUSH PRIVILEGES;