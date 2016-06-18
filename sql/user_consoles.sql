CREATE TABLE `user_consoles` (
  `name` varchar(255) NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `has` tinyint(1) DEFAULT NULL,
  `manual` tinyint(1) DEFAULT NULL,
  `box` tinyint(1) DEFAULT NULL,
  `rating` tinyint(1) DEFAULT NULL,
  `review` varchar(2048) DEFAULT NULL,
  `want` bool default FALSE,
  `wantgames` bool default FALSE,
  PRIMARY KEY (`name`,`user_id`)
)
