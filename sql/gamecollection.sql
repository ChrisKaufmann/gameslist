CREATE TABLE `gamecollection` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `game_id` int(10) unsigned NOT NULL,
  `has_box` tinyint(1) DEFAULT '0',
  `has_manual` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) AUTO_INCREMENT=1;
