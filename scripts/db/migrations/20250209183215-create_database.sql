
-- +migrate Up
CREATE TABLE IF NOT EXISTS `pack_management`.`person` (
  `id` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);
CREATE INDEX `person_name_index` ON `person` (`name`);

CREATE TABLE IF NOT EXISTS `pack_management`.`pack` (
  `id` VARCHAR(255) NOT NULL,
  `description` TEXT NOT NULL,
  `fun_fact` TEXT NULL DEFAULT NULL,
  `is_holiday` BOOLEAN NULL DEFAULT NULL,
  `sender_id` VARCHAR(255) NOT NULL,
  `receiver_id` VARCHAR(255) NOT NULL,
  `status` ENUM('IN_TRANSIT', 'CREATED', 'DELIVERED', 'CANCELED') NOT NULL DEFAULT 'CREATED',
  `estimate_delivery_date` DATE NULL DEFAULT NULL,
  `delivered_at` TIMESTAMP NULL DEFAULT NULL,
  `canceled_at` TIMESTAMP NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`sender_id`) REFERENCES `person`(`id`),
  FOREIGN KEY (`receiver_id`) REFERENCES `person`(`id`)
);

CREATE TABLE IF NOT EXISTS `pack_management`.`pack_event` (
  `id` VARCHAR(255) NOT NULL,
  `pack_id` VARCHAR(255) NOT NULL,
  `location` VARCHAR(255) NOT NULL,
  `description` TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`pack_id`) REFERENCES `pack`(`id`)
);

CREATE TABLE IF NOT EXISTS `pack_management`.`holiday` (
  `id` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `date` DATE NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
);

-- +migrate Down
DROP TABLE `pack_management`.`pack_event`;
DROP TABLE `pack_management`.`pack`;
DROP TABLE `pack_management`.`person`;
DROP TABLE `pack_management`.`holiday`;
