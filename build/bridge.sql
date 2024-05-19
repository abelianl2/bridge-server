/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80036
 Source Host           : localhost:3306
 Source Schema         : bridge

 Target Server Type    : MySQL
 Target Server Version : 80036
 File Encoding         : 65001

 Date: 28/04/2024 22:38:17
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;
CREATE
DATABASE IF NOT EXISTS `bridge`;
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
    `create_time`  datetime DEFAULT CURRENT_TIMESTAMP,
    `update_time`  datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `hash`         varchar(255)                                                 NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash1` (`hash`) USING BTREE,
    KEY            `ids` (`from_network`,`from_address`,`to_network`,`to_address`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
SET
FOREIGN_KEY_CHECKS = 1;
