-- create database
CREATE DATABASE IF NOT EXISTS gigamunch;
-- use database
USE gigamunch;
-- create live_posts table
CREATE TABLE IF NOT EXISTS `live_posts` (
 `post_id` BIGINT NOT NULL PRIMARY KEY,
 `item_id` BIGINT NOT NULL PRIMARY KEY,
 `gigachef_id` VARCHAR(45) NOT NULL,
 `close_datetime` DATETIME NOT NULL, -- used for cron job
 `ready_datetime` DATETIME NOT NULL,
 `search_tags` VARCHAR( 500 ) NOT NULL,
 `is_order_now` BOOLEAN NOT NULL DEFAULT 0,
 `is_experimental` BOOLEAN NOT NULL DEFAULT 0,
 `is_baked_good` BOOLEAN NOT NULL DEFAULT 0,
 `latitude` FLOAT( 10, 6 ) NOT NULL,
 `longitude` FLOAT( 10, 6 ) NOT NULL,
  FULLTEXT(`search_tags`),
  INDEX(`latitude`),
  INDEX(`longitude`),
  INDEX(`ready_datetime`)
) ENGINE = MYISAM ;
-- create like table
CREATE TABLE IF NOT EXISTS `like` (
  `user_id` VARCHAR(45) NOT NULL,
  `item_id` BIGINT NOT NULL,
  `created_datetime` DATETIME NOT NULL,
  PRIMARY KEY (`item_id`, `user_id`)
) ENGINE = InnoDB;
