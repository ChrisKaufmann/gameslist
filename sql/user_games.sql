CREATE TABLE `user_games` (
`id` int(10) unsigned NOT NULL,
`game_id` int(10) unsigned NOT NULL,
`user_id` int(10) unsigned NOT NULL,
`has` tinyint(1) DEFAULT NULL,
`manual` tinyint(1) DEFAULT NULL,
`box` tinyint(1) DEFAULT NULL,
`rating` tinyint(1) DEFAULT NULL,
`review` varchar(2048) DEFAULT NULL,
PRIMARY KEY (`game_id`,`user_id`)
)
