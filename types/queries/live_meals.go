package queries

import "fmt"

var (
	sortByDistance = `SELECT meal_id, ( 3959 * acos( cos( radians(%f) ) * cos( radians( latitude ) ) * cos( radians( longitude ) - radians(%f) ) + sin( radians(%f) ) * sin( radians( latitude ) ) ) ) AS distance
                      FROM live_meals
                      HAVING distance < %d
                      ORDER BY distance
                      LIMIT %d , %d`
	// TODO user prepare statements instead
	searchAndSortByDistance = "SELECT MATCH (search_tags) AGAINST (%s IN NATURAL LANGUAGE MODE)"
)

// GetSortByDistanceQuery returns the query string for getting id, distance from mysql
func GetSortByDistanceQuery(latitude float32, longitude float32, distance int, startLimit int, endLimit int) string {
	return fmt.Sprintf(sortByDistance, latitude, longitude, latitude, distance, startLimit, endLimit)
}
