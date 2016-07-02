 CREATE TABLE `games` (
 `id` int(10) unsigned auto_increment NOT NULL,
 `name` varchar(255) NOT NULL,
 `console_name` varchar(255) DEFAULT NULL,
 `publisher` varchar(255) DEFAULT NULL,
 `year` int(10) unsigned DEFAULT NULL,
 `ebay_price` decimal(5,2) DEFAULT NULL,
 `ebay_updated` datetime DEFAULT NULL,
 `ebay_ends` datetime DEFAULT NULL,
 `ebay_url` varchar(255) DEFAULT NULL,
 PRIMARY KEY (`id`)
 ) ENGINE=InnoDB  AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
