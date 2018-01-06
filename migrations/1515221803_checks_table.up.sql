BEGIN;

ALTER TABLE `proxies` ADD COLUMN `check` int DEFAULT NULL;

CREATE TABLE `checks` (
	`id` int NOT NULL AUTO_INCREMENT,
	`proxy` varchar(21) NOT NULL,
	`state` bool NOT NULL DEFAULT '0',
	`checktime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 DEFAULT COLLATE=utf8_general_ci;

ALTER TABLE `proxies` ADD CONSTRAINT `proxies_fk0` FOREIGN KEY (`check`) REFERENCES `checks`(`id`);

ALTER TABLE `checks` ADD CONSTRAINT `checks_fk0` FOREIGN KEY (`proxy`) REFERENCES `proxies`(`addr`);

COMMIT;
