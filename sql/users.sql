drop table if exists `users`;
create table users (
username varchar(128) primary key not null,
password text(255),
userid varchar(128),
userlevel tinyint(1),
email varchar(128),
timestamp int(11) not null,
token varchar(128)
);
