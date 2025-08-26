-- Active: 1756188174936@@192.168.192.1@3310@kratos_community_user
CREATE TABLE `users`(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_name` VARCHAR(30) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
    `email` VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间（用于软删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_name` (`user_name`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';