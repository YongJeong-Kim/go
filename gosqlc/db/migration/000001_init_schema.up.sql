CREATE TABLE `accounts` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
	`owner` VARCHAR(255) NOT NULL,
	`balance` BIGINT(20) NOT NULL,
	`currency` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL DEFAULT '',
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `entries` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
	`account_id` BIGINT(20) NOT NULL,
	`amount` BIGINT(20) NOT NULL,
	`created_at` DATETIME NOT NULL DEFAULT '',
	PRIMARY KEY (`id`),
	INDEX `FK_entries_accounts` (`account_id`),
	CONSTRAINT `FK_entries_accounts` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `transfers` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
	`from_account_id` BIGINT(20) NOT NULL,
	`to_account_id` BIGINT(20) NOT NULL,
	`amount` BIGINT(20) NOT NULL,
	`created_at` DATETIME NOT NULL DEFAULT '',
	PRIMARY KEY (`id`),
	INDEX `FK_transfers_accounts` (`from_account_id`),
	INDEX `FK_transfers_accounts_2` (`to_account_id`),
	CONSTRAINT `FK_transfers_accounts` FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`),
	CONSTRAINT `FK_transfers_accounts_2` FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
