-- Active: 1756839858382@@192.168.192.1@3310@kratos_community_interaction
CREATE TABLE `likes` (
  `user_id` bigint unsigned NOT NULL COMMENT '点赞的用户ID',
  `article_id` bigint unsigned NOT NULL COMMENT '被点赞的文章ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
  -- 关键：使用 user_id 和 article_id 组成复合主键，确保唯一性
  PRIMARY KEY (`user_id`, `article_id`),
  -- 为了反向查询（查一篇文章的所有点赞），也为 article_id 创建一个索引
  KEY `idx_article_id` (`article_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞关系表';