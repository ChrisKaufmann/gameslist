CREATE TABLE `ratings` (
	`thing_id` int(10) unsigned NOT NULL,
	`user_id` int(10) unsigned NOT NULL,
	`rating` int(10) unsigned DEFAULT '0',
	PRIMARY KEY (`thing_id`)
);
