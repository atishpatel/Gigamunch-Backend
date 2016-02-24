CREATE TABLE IF NOT EXISTS `live_meals` (
 `meal_id` BIGINT NOT NULL PRIMARY KEY,
 `close_datetime` DATETIME NOT NULL, -- user for cron job
 `ready_datetime` DATETIME NOT NULL,
 `search_tags` VARCHAR( 500 ) NOT NULL,
 `is_experimental` BOOLEAN NOT NULL DEFAULT 0,
 `is_baked_good` BOOLEAN NOT NULL DEFAULT 0,
 `latitude` FLOAT( 10, 6 ) NOT NULL,
 `longitude` FLOAT( 10, 6 ) NOT NULL,
  FULLTEXT(`search_tags`),
  INDEX(`latitude`),
  INDEX(`longitude`)
) ENGINE = MYISAM ;
