-- Active: 1756839858382@@192.168.192.1@3310@kratos_community_relation
CREATE TABLE `relations` (
  `follower_id` bigint unsigned NOT NULL COMMENT '粉丝ID (发起关注的用户)',
  `following_id` bigint unsigned NOT NULL COMMENT '被关注者ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '关注时间',
  PRIMARY KEY (`follower_id`,`following_id`),
  KEY `idx_following_id` (`following_id`) COMMENT '被关注者ID索引，用于快速查找粉丝列表'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户关注关系表';