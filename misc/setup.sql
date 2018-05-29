-- create database
CREATE DATABASE IF NOT EXISTS gigamunch CHARACTER SET utf8mb4;
-- use database
USE gigamunch;
-- create live_posts table
-- CREATE TABLE IF NOT EXISTS `active_items` (
--   `id` BIGINT NOT NULL PRIMARY KEY,
--   `menu_id` BIGINT NOT NULL,
--   `cook_id` VARCHAR(45) NOT NULL,
--   `created_datetime` DATETIME NOT NULL,
--   `cook_price_per_serving` FLOAT(10, 2) NOT NULL,
--   `min_servings` TINYINT UNSIGNED NOT NULL,
--   `max_servings` SMALLINT UNSIGNED NOT NULL,
--   `latitude` FLOAT( 10, 6 ) NOT NULL,
--   `longitude` FLOAT( 10, 6 ) NOT NULL,
--   `vegan` BOOLEAN NOT NULL DEFAULT 0,
--   `vegetarian` BOOLEAN NOT NULL DEFAULT 0, 
--   `paleo` BOOLEAN NOT NULL DEFAULT 0, 
--   `gluten_free` BOOLEAN NOT NULL DEFAULT 0, 
--   `kosher` BOOLEAN NOT NULL DEFAULT 0,
--   INDEX(`latitude`),
--   INDEX(`longitude`),
--   INDEX(`created_datetime`),
--   INDEX(`menu_id`)
-- ) ENGINE = InnoDB CHARACTER SET utf8mb4;
-- create like table
-- CREATE TABLE IF NOT EXISTS `likes` (
--   `user_id` VARCHAR(45) NOT NULL,
--   `item_id` BIGINT NOT NULL,
--   `cook_id` VARCHAR(45) NOT NULL,
--   `menu_id` BIGINT NOT NULL,
--   `created_datetime` DATETIME NOT NULL,
--   PRIMARY KEY (`item_id`, `user_id`)
-- ) ENGINE = InnoDB;
-- create review table
-- CREATE TABLE IF NOT EXISTS `review` (
--   `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
--   `cook_id` VARCHAR(45) NOT NULL,
--   `eater_id` VARCHAR(45) NOT NULL,
--   `eater_name` VARCHAR(100) NOT NULL,
--   `eater_photo_url` VARCHAR(350) NOT NULL,
--   `inquiry_id` BIGINT NOT NULL,
--   `item_id` BIGINT NOT NULL,
--   `item_photo_url` VARCHAR(350) NOT NULL,
--   `item_name` VARCHAR(100) NOT NULL,
--   `menu_id` BIGINT NOT NULL,
--   `created_datetime` DATETIME NOT NULL DEFAULT NOW(),
--   `rating` TINYINT NOT NULL,
--   `text` VARCHAR(1200),
--   `edited_datetime` DATETIME NOT NULL DEFAULT NOW(),
--   `is_edited` BOOLEAN NOT NULL DEFAULT 0,
--   -- cook response stuff
--   `has_response` BOOLEAN NOT NULL DEFAULT 0,
--   `response_created_datetime` DATETIME,
--   `response_text` VARCHAR(1200),
--   INDEX(`item_id`),
--   INDEX(`created_datetime`)
-- ) ENGINE = InnoDB CHARACTER SET utf8mb4;
-- create promo_codes table
CREATE TABLE IF NOT EXISTS `promo_code` (
  `code` VARCHAR(45) NOT NULL PRIMARY KEY,
  `created_datetime` DATETIME NOT NULL DEFAULT NOW(),
  `free_delivery` BOOLEAN NOT NULL DEFAULT 0,
  `percent_off` TINYINT NOT NULL DEFAULT 0,
  `amount_off` FLOAT( 6, 2 ) NOT NULL DEFAULT 0,
  `discount_cap` FLOAT(6,2) NOT NULL DEFAULT 0,
  `free_dish` BOOLEAN NOT NULL DEFAULT 0,
  `buy_one_get_one_free` BOOLEAN NOT NULL DEFAULT 0,
  `start_datetime` DATETIME NOT NULL DEFAULT NOW(),
  `end_datetime` DATETIME NOT NULL DEFAULT NOW(),
  `num_uses` INT NOT NULL DEFAULT 0
) ENGINE = InnoDB CHARACTER SET utf8mb4;
-- create used_promo_code table 
-- CREATE TABLE IF NOT EXISTS `used_promo_code` (
--   `eater_id` VARCHAR(45) NOT NULL,
--   `eater_name` VARCHAR(100) NOT NULL,
--   `inquiry_id` BIGINT NOT NULL,
--   `created_datetime` DATETIME NOT NULL DEFAULT NOW(),
--   `code` VARCHAR(45) NOT NULL,
--   `state` TINYINT NOT NULL DEFAULT 0,
--   PRIMARY KEY (`eater_id`, `inquiry_id`)
-- ) ENGINE = InnoDB CHARACTER SET utf8mb4;
-- create completed_inquries table 
-- CREATE TABLE IF NOT EXISTS `completed_inquiries` (
--   `id` BIGINT NOT NULL PRIMARY KEY,
--   `servings` SMALLINT UNSIGNED NOT NULL,
--   `cook_id` VARCHAR(45) NOT NULL,
--   `item_id` BIGINT NOT NULL,
--   `menu_id` BIGINT NOT NULL,
--   `created_datetime` DATETIME NOT NULL,
--   `cook_price_per_serving` FLOAT(10,2) NOT NULL,
--   `cook_revenue` FLOAT(10,2) NOT NULL,
--   INDEX(`cook_id`)
-- ) ENGINE = InnoDB;
-- create sub
CREATE TABLE IF NOT EXISTS `sub`(
  `date` DATE NOT NULL,
  `sub_email` VARCHAR(175) NOT NULL,
  `created_datetime` DATETIME NOT NULL DEFAULT NOW(),
  `skip` BOOLEAN NOT NULL DEFAULT 0,
  `servings` TINYINT NOT NULL,
  `amount` FLOAT(6,2) NOT NULL,
  `amount_paid` FLOAT(6,2) NOT NULL DEFAULT 0,
  `paid` BOOLEAN NOT NULL DEFAULT 0,
  `paid_datetime` DATETIME,
  `refunded` BOOLEAN NOT NULL DEFAULT 0,
  `delivery_time` TINYINT NOT NULL,
  `payment_method_token` VARCHAR(37) NOT NULL DEFAULT '',
  `transaction_id` VARCHAR(37) NOT NULL DEFAULT '',
  `free` BOOLEAN NOT NULL DEFAULT 0,
  `discount_amount` FLOAT(6,2) NOT NULL DEFAULT 0,
  `discount_percent` TINYINT NOT NULL DEFAULT 0,
  `customer_id` VARCHAR(37) NOT NULL,
   PRIMARY KEY (`date`, `sub_email`)
) ENGINE = InnoDB;
-- create coupon
-- create promo
-- create deliveries
CREATE TABLE IF NOT EXISTS deliveries (
	date DATE NOT NULL,
	sub_email VARCHAR(175) NOT NULL,
	created_dt  DATETIME NOT NULL DEFAULT NOW(),
	updated_dt DATETIME NOT NULL DEFAULT NOW(),
	driver_id BIGINT NOT NULL DEFAULT -1,
	driver_email VARCHAR(175) NOT NULL,
	driver_name VARCHAR(125) NOT NULL,
	sub_id BIGINT NOT NULL DEFAULT -1,
	delivery_order SMALLINT NOT NULL DEFAULT -1,
	success BOOLEAN DEFAULT 0,
	fail BOOLEAN DEFAULT 0,
	PRIMARY KEY (date, sub_email)
)
INSERT INTO deliveries (driver_email, driver_name, sub_email,date) VALUES ('d1', 'd1n', 'sub1', '2018-06-28')
-- create activity
CREATE TABLE IF NOT EXISTS activity (
	created_dt DATETIME NOT NULL DEFAULT NOW(),
	date DATE NOT NULL,
	user_id BIGINT NOT NULL,
	email VARCHAR(175) NOT NULL,
	first_name VARCHAR(125) NOT NULL,
	last_name VARCHAR(125) NOT NULL,
	location INT NOT NULL,
	addr_changed BOOLEAN NOT NULL DEFAULT 0,
	addr_apt VARCHAR(50) NOT NULL DEFAULT '',
	addr_string VARCHAR(175) NOT NULL,
	zip VARCHAR(30) NOT NULL,
	lat FLOAT( 10, 6 ) NOT NULL,
	`long` FLOAT( 10, 6 ) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT 1,
	skip BOOLEAN NOT NULL DEFAULT 0,
	veg_bag BOOLEAN NOT NULL DEFAULT 0,
	servings TINYINT NOT NULL,
	servings_changed BOOLEAN NOT NULL DEFAULT 0,
	first BOOLEAN NOT NULL DEFAULT 0,
	amount FLOAT(6,2) NOT NULL,
	amount_paid FLOAT(6,2) NOT NULL DEFAULT 0,
	discount_amount FLOAT(6,2) NOT NULL DEFAULT 0,
	paid BOOLEAN NOT NULL DEFAULT 0,
	paid_dt DATETIME,
	transaction_id VARCHAR(37) NOT NULL DEFAULT '',
	refunded BOOLEAN NOT NULL DEFAULT 0,
	refunded_amount FLOAT(6,2) NOT NULL DEFAULT 0,
	refunded_dt DATETIME,
	refund_transaction_id VARCHAR(37) NOT NULL DEFAULT '',
	coupon_id BIGINT,
	payment_provider TINYINT NOT NULL DEFAULT 0,
	forgiven BOOLEAN NOT NULL DEFAULT 0,
	gift BOOLEAN NOT NULL DEFAULT 0,
	gift_from_user_id BIGINT,
	deviant BOOLEAN NOT NULL DEFAULT 0,
	deviant_reason VARCHAR(225) NOT NULL DEFAULT '',
	PRIMARY KEY (date, user_id)
) ENGINE = InnoDB CHARACTER SET utf8mb4;