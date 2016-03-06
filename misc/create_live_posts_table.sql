CREATE TABLE IF NOT EXISTS `live_posts` (
 `post_id` BIGINT NOT NULL PRIMARY KEY,
 `close_datetime` DATETIME NOT NULL, -- used for cron job
 `ready_datetime` DATETIME NOT NULL,
 `search_tags` VARCHAR( 500 ) NOT NULL,
 `is_experimental` BOOLEAN NOT NULL DEFAULT 0,
 `is_baked_good` BOOLEAN NOT NULL DEFAULT 0,
 `latitude` FLOAT( 10, 6 ) NOT NULL,
 `longitude` FLOAT( 10, 6 ) NOT NULL,
  FULLTEXT(`search_tags`),
  INDEX(`latitude`),
  INDEX(`longitude`),
  INDEX(`ready_datetime`)
) ENGINE = MYISAM ;
