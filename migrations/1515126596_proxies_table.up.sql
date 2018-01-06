BEGIN;

DROP TABLE IF EXISTS `proxies`;
CREATE TABLE `proxies` (
	`addr` varchar(21) NOT NULL,
	`class` tinyint(1) NOT NULL DEFAULT '0',
	`anon` tinyint(1) NOT NULL DEFAULT '0',
	`created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`addr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 DEFAULT COLLATE=utf8_general_ci;

COMMIT;
