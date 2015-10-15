 CREATE TABLE `things` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `parent_id` int(10) unsigned DEFAULT NULL,
  `type` enum('console','game','manual','box') DEFAULT NULL,
  `rating` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
)
