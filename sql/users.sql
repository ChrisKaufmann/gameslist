drop table if exists `users`;
create table users (
id int unsigned primary key not null auto_increment,
email varchar(128),
admin bool default false,
share_token char(128),
login_token char(128)
);
