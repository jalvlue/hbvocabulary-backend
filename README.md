# 数据库
- mysql:8.3.0
- `CREATE USER 'user_course_design'@'%' IDENTIFIED BY '123456';`, 新建用户: user_course_design, 密码: 123456
- `GRANT CREATE ON *.* TO 'user_course_design'@'%';`, 给用户授权
- `FLUSH PRIVILEGES;`, 刷新 mysql 权限

- 建库建表：
```sql
create database db_course_design;

CREATE TABLE `t_users` (
  `id` bigserial PRIMARY KEY,
  `username` varchar(255) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `test_count` int NOT NULL DEFAULT 0,
  `max_score` int NOT NULL DEFAULT 0,
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE `t_vocabularies` (
  `id` bigserial PRIMARY KEY,
  `word` varchar(255) NOT NULL,
  `difficulty` smallint NOT NULL
);

CREATE TABLE `t_test_records` (
  `id` bigserial PRIMARY KEY,
  `user_id` bigint NOT NULL,
  `score` int NOT NULL,
  `created_at` timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE
  `t_test_records`
ADD
  FOREIGN KEY (`user_id`) REFERENCES `t_users` (`id`);
```