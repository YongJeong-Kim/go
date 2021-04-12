CREATE TABLE `users` (
  `username` VARCHAR(255) NOT NULL,
  `hashed_password` VARCHAR(255) NOT NULL,
  `full_name` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `password_changed_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(username)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

ALTER TABLE `accounts` ADD CONSTRAINT `fk_accounts_owner` FOREIGN KEY(`owner`) REFERENCES `users` (`username`);

-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE `accounts` ADD CONSTRAINT `owner_currency_key` UNIQUE (`owner`, `currency`);