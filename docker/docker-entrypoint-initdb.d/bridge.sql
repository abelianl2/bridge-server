CREATE
DATABASE IF NOT EXISTS `bridge`;
USE bridge;
SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;
-- ----------------------------
-- Table structure for bridge
-- ----------------------------
DROP TABLE IF EXISTS `deposit`;
CREATE TABLE `deposit`
(
    `id`           bigint                                                       NOT NULL AUTO_INCREMENT,
    `from_network` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `from_address` varchar(255)                                                 NOT NULL,
    `to_network`   varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `to_address`   varchar(255)                                                 NOT NULL,
    `call_back`    text,
    `create_time`  datetime                                                      DEFAULT CURRENT_TIMESTAMP,
    `update_time`  datetime                                                      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `hash`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
    `uuid`         varchar(255)                                                 NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uuid` (`uuid`) USING BTREE,
    KEY            `ids` (`from_network`,`from_address`,`to_network`,`to_address`) USING BTREE,
    KEY            `hash1` (`hash`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
SET
FOREIGN_KEY_CHECKS = 1;
