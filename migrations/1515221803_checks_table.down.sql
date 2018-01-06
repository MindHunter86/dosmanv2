BEGIN;

ALTER TABLE `checks` DROP FOREIGN KEY `checks_fk0`;
ALTER TABLE `proxies` DROP FOREIGN KEY `proxies_fk0`;

DROP TABLE checks;

ALTER TABLE `proxies` DROP COLUMN `check`;

COMMIT;
