-- Active: 1756839736469@@192.168.192.1@3310@kratos_community_user
CREATE TABLE `articles`(
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文章ID',
    `author_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '作者的用户ID',
    `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文章标题',
    `content` longtext COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文章内容',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_author_id` (`author_id`)
)