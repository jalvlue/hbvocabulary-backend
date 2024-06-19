# 数据库
- mysql:8.3.0
- `CREATE USER 'course_design_user'@'%' IDENTIFIED BY '123456';`, 新建用户: course_design_user, 密码: 123456
- `GRANT CREATE ON *.* TO 'course_design_user'@'%';`, 给用户授权
- `FLUSH PRIVILEGES;`, 刷新 mysql 权限

- `create database db_course_design;`, 新建数据库
