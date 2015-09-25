drop table if exists `users`;
create table users (
id int unsigned primary key not null auto_increment,
email varchar(128)
);
