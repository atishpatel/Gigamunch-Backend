CREATE TABLE `live_meals` (
 `meal_id` BIGINT NOT NULL PRIMARY KEY,
 `search_tags` VARCHAR( 300 ) NOT NULL ,
 `latitude` FLOAT( 10, 6 ) NOT NULL ,
 `longitude` FLOAT( 10, 6 ) NOT NULL ,
  FULLTEXT(`search_tags`),
  INDEX(`latitude`),
  INDEX(`longitude`)
) ENGINE = MYISAM ;
