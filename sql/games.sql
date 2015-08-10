drop table if exists games;
CREATE TABLE `games` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `console_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`)
);
