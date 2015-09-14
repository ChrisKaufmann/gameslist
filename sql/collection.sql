CREATE TABLE `collection` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `thing_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`)
);
