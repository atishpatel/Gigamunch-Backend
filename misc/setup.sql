-- create database
CREATE DATABASE IF NOT EXISTS gigamunch CHARACTER SET utf8mb4;
-- use database
USE gigamunch;
-- create live_posts table
CREATE TABLE IF NOT EXISTS `active_items` (
  `id` BIGINT NOT NULL PRIMARY KEY,
  `menu_id` BIGINT NOT NULL,
  `cook_id` VARCHAR(45) NOT NULL,
  `created_datetime` DATETIME NOT NULL,
  `cook_price_per_serving` FLOAT(10, 2) NOT NULL,
  `min_servings` TINYINT UNSIGNED NOT NULL,
  `max_servings` SMALLINT UNSIGNED NOT NULL,
  `latitude` FLOAT( 10, 6 ) NOT NULL,
  `longitude` FLOAT( 10, 6 ) NOT NULL,
  `vegan` BOOLEAN NOT NULL DEFAULT 0,
  `vegetarian` BOOLEAN NOT NULL DEFAULT 0, 
  `paleo` BOOLEAN NOT NULL DEFAULT 0, 
  `gluten_free` BOOLEAN NOT NULL DEFAULT 0, 
  `kosher` BOOLEAN NOT NULL DEFAULT 0,
  INDEX(`latitude`),
  INDEX(`longitude`),
  INDEX(`created_datetime`),
  INDEX(`menu_id`)
) ENGINE = MYISAM;
-- create like table
CREATE TABLE IF NOT EXISTS `likes` (
  `user_id` VARCHAR(45) NOT NULL,
  `item_id` BIGINT NOT NULL,
  `cook_id` VARCHAR(45) NOT NULL,
  `menu_id` BIGINT NOT NULL,
  `created_datetime` DATETIME NOT NULL,
  PRIMARY KEY (`item_id`, `user_id`)
) ENGINE = InnoDB;
-- CREATE TABLE IF NOT EXISTS `like` (
--   `user_id` VARCHAR(45) NOT NULL,
--   `item_id` BIGINT NOT NULL,
--   `created_datetime` DATETIME NOT NULL,
--   PRIMARY KEY (`item_id`, `user_id`)
-- ) ENGINE = InnoDB;
-- create num_likes table
-- CREATE TABLE IF NOT EXISTS `num_likes` (
--   `item_id` BIGINT NOT NULL PRIMARY KEY,
--   `num` BIGINT NOT NULL
-- ) ENGINE = InnoDB;
-- create review table
CREATE TABLE IF NOT EXISTS `review` (
  `id` BIGINT NOT NULL PRIMARY KEY AUTOINCREMENT,
  `cook_id` VARCHAR(45) NOT NULL,
  `eater_id` VARCHAR(45) NOT NULL,
  `eater_name` VARCHAR(100) NOT NULL,
  `eater_photo_url` VARCHAR(200) NOT NULL,
  `inquiry_id` BIGINT NOT NULL,
  `item_id` BIGINT NOT NULL,
  `menu_id` BIGINT NOT NULL,
  `created_datetime` DATETIME NOT NULL DEFAULT NOW(),
  `rating` INT NOT NULL,
  `text` VARCHAR(1200),
  `edited_datetime` DATETIME NOT NULL DEFAULT NOW(),
  `is_edited` BOOLEAN NOT NULL DEFAULT 0,
  -- cook response stuff
  `has_response` BOOLEAN NOT NULL DEFAULT 0,
  `response_created_datetime` DATETIME,
  `response_text` VARCHAR(1200),
  INDEX(`item_id`),
  INDEX(`created_datetime`)
) ENGINE = InnoDB;
-- create completed_inquries table 
CREATE TABLE IF NOT EXISTS `completed_inquries` (
  `id` BIGINT NOT NULL PRIMARY KEY,
  `servings` SMALLINT UNSIGNED NOT NULL,
  `cook_id` VARCHAR(45) NOT NULL,
  `item_id` BIGINT NOT NULL,
  `menu_id` BIGINT NOT NULL,
  `created_datetime` DATETIME NOT NULL,
  `cook_price_per_serving` FLOAT(10,2) NOT NULL,
  `cook_revenue` FLOAT(10,2) NOT NULL,
  INDEX(`cook_id`)
) ENGINE = InnoDB;
