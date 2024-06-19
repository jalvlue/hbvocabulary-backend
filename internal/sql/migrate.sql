CREATE USER 'user_course_design' @'%' IDENTIFIED BY '123456';

GRANT CREATE ON *.* TO 'course_design_user' @'%';

FLUSH PRIVILEGES;

create database db_course_design;