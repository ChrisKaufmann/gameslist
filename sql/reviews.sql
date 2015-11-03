DROP TABLE IF EXISTS `reviews`;
CREATE TABLE `reviews` (
	`thing_id` int(10) unsigned NOT NULL,
	`user_id` int(10) unsigned NOT NULL,
	`review` text
);
