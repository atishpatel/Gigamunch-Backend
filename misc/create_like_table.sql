CREATE TABLE `like` (
  `user_id` VARCHAR(45) NOT NULL,
  `item_id` BIGINT NOT NULL,
  `created_datetime` DATETIME NOT NULL,
  PRIMARY KEY (`item_id`, `user_id`)
) ENGINE = InnoDB;
