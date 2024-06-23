-- CREATE USER 'user_course_design' @'%' IDENTIFIED BY '123456';
-- GRANT ALL ON *.* TO 'user_course_design' @'%';
-- FLUSH PRIVILEGES;
create database db_course_design;

use db_course_design;

CREATE TABLE `t_users` (
  `id` int PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(255) UNIQUE NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `test_count` int NOT NULL DEFAULT 0 COMMENT '测试次数',
  `max_score` int NOT NULL DEFAULT 0 COMMENT '最高分数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);

CREATE TABLE `t_vocabularies` (
  `id` int PRIMARY KEY AUTO_INCREMENT COMMENT '词汇ID',
  `word` varchar(255) NOT NULL COMMENT '单词' -- `difficulty` smallint NOT NULL COMMENT '难度'
);

CREATE TABLE `t_test_records` (
  `id` int PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
  `user_id` int NOT NULL COMMENT '用户ID',
  `score` int NOT NULL COMMENT '分数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);

ALTER TABLE
  `t_test_records`
ADD
  FOREIGN KEY (`user_id`) REFERENCES `t_users` (`id`);