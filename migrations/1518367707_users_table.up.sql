BEGIN;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
	`steamid64` int(17) UNSIGNED NOT NULL,
	`partner_id` int(9) UNSIGNED NOT NULL UNIQUE,
	`trade_token` varchar(8) NOT NULL,
	`remember_token` varchar(64) NOT NULL,
	`username` varchar(64) DEFAULT NULL,
	`logged_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`steamid64`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 DEFAULT COLLATE=utf8_general_ci;

COMMIT;
