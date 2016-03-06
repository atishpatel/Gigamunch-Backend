SELECT item_id, ( 3959 * acos( cos( radians(35.1) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(-85.1) ) + sin( radians(35.1) ) * sin( radians( latitude ) ) ) ) AS distance
FROM live_posts
WHERE ready_datetime
BETWEEN '2016-01-03 01:00:00' AND '3000-01-01 00:00:00'
HAVING distance < 50
ORDER BY ready_datetime ASC, distance
LIMIT 0 , 10;

SELECT item_id, ( 3959 * acos( cos( radians(35.1) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(-85.1) ) + sin( radians(35.1) ) * sin( radians( latitude ) ) ) ) AS distance
FROM live_posts
WHERE ready_datetime
BETWEEN '2000-01-01 00:00:00' AND '2016-01-03 01:00:00'
HAVING distance < 50
ORDER BY ready_datetime ASC, distance
LIMIT 0 , 10;
