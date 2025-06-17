-- 创建数据库
CREATE DATABASE IF NOT EXISTS `knowledge_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE `knowledge_system`;

-- 创建知识表
CREATE TABLE IF NOT EXISTS `knowledge` (
  `id` varchar(36) NOT NULL COMMENT '唯一ID，服务端生成uuid',
  `content` text NOT NULL COMMENT '知识内容',
  `labels` json DEFAULT NULL COMMENT '标签分数数组，存储为JSON字符串',
  `summary` text DEFAULT NULL COMMENT '内容摘要',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  FULLTEXT KEY `idx_content` (`content`) COMMENT '内容全文索引',
  FULLTEXT KEY `idx_summary` (`summary`) COMMENT '摘要全文索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识条目表';

-- 创建用户并授权（根据实际情况修改用户名和密码）
-- CREATE USER 'knowledge_user'@'%' IDENTIFIED BY 'knowledge_password';
-- GRANT ALL PRIVILEGES ON knowledge_system.* TO 'knowledge_user'@'%';
-- FLUSH PRIVILEGES; 