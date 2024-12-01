CREATE TABLE `account`.`account`  (
                                       `id` int UNSIGNED NOT NULL AUTO_INCREMENT,
                                       `account_id` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
                                       `account_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
                                       `password` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
                                       `created_at` bigint UNSIGNED NOT NULL DEFAULT 0,
                                       `updated_at` bigint UNSIGNED NOT NULL DEFAULT 0,
                                       PRIMARY KEY (`id`) USING BTREE,
                                       UNIQUE INDEX `account_id_unique`(`account_id` ASC) USING BTREE,
                                       UNIQUE INDEX `account_name_unique`(`account_name` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;