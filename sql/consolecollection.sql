CREATE TABLE `consolecollection` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`user_id` int(10) unsigned NOT NULL,
	`console_id` int(10) unsigned NOT NULL,
	`has_box` bool default false,
	`has_manual` bool default false,
	PRIMARY KEY (`id`)
) AUTO_INCREMENT=1
