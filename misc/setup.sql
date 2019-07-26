-- create database
CREATE DATABASE IF NOT EXISTS gigamunch CHARACTER SET utf8mb4;
-- use database
USE gigamunch;
-- create driver_assignment
CREATE TABLE IF NOT EXISTS deliveries (
	created_dt DATETIME NOT NULL DEFAULT NOW(),
	date DATE NOT NULL,
	driver_id VARCHAR(125) NOT NULL,
	driver_email VARCHAR(175) NOT NULL,
	driver_name VARCHAR(125) NOT NULL,
	sub_id VARCHAR(125) NOT NULL,
	sub_email VARCHAR(175) NOT NULL,
	order INT NOT NULL DEFAULT -1,
	delivered BOOLEAN DEFAULT 0,
	PRIMARY KEY (date, sub_email)
);
-- create activity
CREATE TABLE IF NOT EXISTS activity(
	created_dt DATETIME NOT NULL DEFAULT NOW(),
	date DATE NOT NULL,
	user_id VARCHAR(125) NOT NULL DEFAULT '',
	email VARCHAR(175) NOT NULL,
	first_name VARCHAR(125) NOT NULL DEFAULT '',
	last_name VARCHAR(125) NOT NULL DEFAULT '',
	location INT NOT NULL DEFAULT 0,
	addr_changed BOOLEAN NOT NULL DEFAULT 0,
	addr_apt VARCHAR(50) NOT NULL DEFAULT '',
	addr_string VARCHAR(175) NOT NULL DEFAULT '',
	zip VARCHAR(30) NOT NULL DEFAULT 0,
	lat FLOAT( 10, 6 ) NOT NULL DEFAULT 0,
	`long` FLOAT( 10, 6 ) NOT NULL DEFAULT 0,
	active BOOLEAN NOT NULL DEFAULT 1,
	skip BOOLEAN NOT NULL DEFAULT 0,
	servings TINYINT NOT NULL DEFAULT 0,
	veg_servings TINYINT NOT NULL DEFAULT 0,
	servings_changed BOOLEAN NOT NULL DEFAULT 0,
	first BOOLEAN NOT NULL DEFAULT 0,
	amount FLOAT(6,2) NOT NULL,
	amount_paid FLOAT(6,2) NOT NULL DEFAULT 0,
	discount_amount FLOAT(6,2) NOT NULL DEFAULT 0,
	discount_percent TINYINT NOT NULL DEFAULT 0,
	paid BOOLEAN NOT NULL DEFAULT 0,
	paid_dt DATETIME,
	transaction_id VARCHAR(37) NOT NULL DEFAULT '',
	payment_method_token VARCHAR(37) NOT NULL DEFAULT '',
	customer_id VARCHAR(37) NOT NULL DEFAULT '',
	refunded BOOLEAN NOT NULL DEFAULT 0,
	refunded_amount FLOAT(6,2) NOT NULL DEFAULT 0,
	refunded_dt DATETIME,
	refund_transaction_id VARCHAR(37) NOT NULL DEFAULT '',
	payment_provider TINYINT NOT NULL DEFAULT 0,
	forgiven BOOLEAN NOT NULL DEFAULT 0,
	gift BOOLEAN NOT NULL DEFAULT 0,
	gift_from_user_id BIGINT,
	deviant BOOLEAN NOT NULL DEFAULT 0,
	deviant_reason VARCHAR(225) NOT NULL DEFAULT '',
	PRIMARY KEY (date, user_id)
) ENGINE = InnoDB CHARACTER SET utf8mb4;

-- TODO: change primary key from date, email to date, user_id
-- create discount
CREATE TABLE IF NOT EXISTS discount(
	id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
	created_dt DATETIME NOT NULL DEFAULT NOW(),
	user_id VARCHAR(125) NOT NULL DEFAULT '',
	email VARCHAR(175) NOT NULL,
	first_name VARCHAR(125) NOT NULL DEFAULT '',
	last_name VARCHAR(125) NOT NULL DEFAULT '',
	date_used DATE NOT NULL DEFAULT '0000-00-00',
	discount_amount FLOAT(6,2) NOT NULL DEFAULT 0,
	discount_percent TINYINT NOT NULL DEFAULT 0
) ENGINE = InnoDB CHARACTER SET utf8mb4;
