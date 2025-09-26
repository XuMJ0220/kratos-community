-- Active: 1756188174936@@192.168.192.1@3310@kratos_community_content
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

CREATE TABLE `outbox_messages`(
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `topic` VARCHAR(255) NOT NULL COMMENT 'kafka topic',
    `message_key` VARCHAR(255) NOT NULL COMMENT '消息 Key',
    `message_value` VARCHAR(255) NOT NULL COMMENT '消息 Value',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态: 0-待发送, 1-已发送',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_status_id` (`status`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='事务性发件箱表';