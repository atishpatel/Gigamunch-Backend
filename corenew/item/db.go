package item

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const (
	datetimeFormat  = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"
	insertStatement = "INSERT INTO `active_items` (id, menu_id, cook_id, created_datetime, cook_price_per_serving, min_servings, max_servings, latitude, longitude, vegan, vegetarian, paleo, gluten_free, kosher) VALUES (%d, %d, '%s', '%s', %f, %d, %d, %f, %f, %t, %t, %t, %t, %t)"
	updateStatement = "UPDATE `active_items` SET menu_id=%d AND cook_price_per_serving=%f AND min_servings=%d AND max_servings=%d AND latitude=%f AND longitude=%f AND vegan=%t AND vegetarian=%t AND paleo=%t AND gluten_free=%t AND kosher=%t WHERE id=%d"
	deleteStatement = "DELETE FROM active_items WHERE id=%d"
	// This select statement gets items sorted by distance and newness. A logarithmic function is applied to newness to the magnititude that an item 85 miles away is ranked higher than an item newly created (a month old item drops to the same rank as one 7 miles away)
	selectStatement = "SELECT a.id, a.menu_id, a.cook_id FROM active_items as a, (SELECT menu_id,(3959*acos(cos(radians(%f))*cos(radians(latitude))*cos(radians(longitude)-radians(%f))+sin(radians(%f))*sin(radians(latitude))))+LEAST(log((Now()-MAX(created_datetime))/200000000)*10,0) as s FROM active_items GROUP BY menu_id ORDER BY s LIMIT %d, %d) as b WHERE b.menu_id = a.menu_id"
)

var (
	connectOnce = sync.Once{}
	mysqlDB     *sql.DB
)

func get(ctx context.Context, id int64) (*Item, error) {
	item := new(Item)
	key := datastore.NewKey(ctx, kindItem, "", id, nil)
	err := datastore.Get(ctx, key, item)
	item.ID = id
	return item, err
}

func put(ctx context.Context, id int64, item *Item) error {
	var err error
	item.ID = id
	key := datastore.NewKey(ctx, kindItem, "", id, nil)
	_, err = datastore.Put(ctx, key, item)
	return err
}

func putIncomplete(ctx context.Context, item *Item) (int64, error) {
	var err error
	key := datastore.NewIncompleteKey(ctx, kindItem, nil)
	key, err = datastore.Put(ctx, key, item)
	item.ID = key.IntID()
	return key.IntID(), err
}

func getRankedActiveItem(ctx context.Context, lat, long float64, startLimit, endLimit int) ([]int64, []int64, []string, error) {
	var err error
	st := fmt.Sprintf(selectStatement, lat, long, lat, startLimit, endLimit)
	rows, err := mysqlDB.Query(st)
	if err != nil {
		return nil, nil, nil, errSQLDB.WithError(err).Wrapf("failed to run query: %s", st)
	}
	defer handleCloser(ctx, rows)
	var itemIDs []int64
	var itemID int64
	var menuIDs []int64
	var menuID int64
	var cookIDs []string
	var cookID string
	for rows.Next() {
		_ = rows.Scan(&itemID, &menuID, &cookID)
		itemIDs = append(itemIDs, itemID)
		menuIDs = append(menuIDs, menuID)
		cookIDs = append(cookIDs, cookID)
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, nil, errSQLDB.WithError(err).Wrap("failed at scanning rows")
	}
	return itemIDs, menuIDs, cookIDs, nil
}

// getCookItems returns a list of Items by ordered by MenuID
func getCookItems(ctx context.Context, cookID string) ([]int64, []Item, error) {
	query := datastore.NewQuery(kindItem).
		Filter("CookID =", cookID).
		Order("MenuID").
		Limit(1000)
	var results []Item
	keys, err := query.GetAll(ctx, &results)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int64, len(keys))
	for i := range keys {
		ids[i] = keys[i].IntID()
	}
	return ids, results, nil
}

// getMulti gets a list of Items
func getMulti(ctx context.Context, ids []int64) ([]Item, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids cannot be 0 for getMulti")
	}
	dst := make([]Item, len(ids))
	var err error
	keys := make([]*datastore.Key, len(ids))
	for i := range ids {
		keys[i] = datastore.NewKey(ctx, kindItem, "", ids[i], nil)
	}
	err = datastore.GetMulti(ctx, keys, dst)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return nil, err
	}
	return dst, nil
}

func insertOrUpdateActiveItem(ctx context.Context, id int64, item *Item, lat, long float64) error {
	if item.Active {
		st := fmt.Sprintf(updateStatement, item.MenuID, item.CookPricePerServing, item.MinServings, item.MaxServings, lat, long,
			item.DietaryConcerns.vegan(), item.DietaryConcerns.vegetarian(), item.DietaryConcerns.paleo(), item.DietaryConcerns.glutenFree(), item.DietaryConcerns.kosher(),
			id)
		_, err := mysqlDB.Exec(st)
		if err != nil {
			return errSQLDB.WithError(err).Wrapf("failed to execute: %s", st)
		}
		// update
		return nil
	}
	// insert
	st := fmt.Sprintf(insertStatement, id, item.MenuID, item.CookID,
		item.CreatedDateTime.UTC().Format(datetimeFormat), item.CookPricePerServing,
		item.MinServings, item.MaxServings, lat, long,
		item.DietaryConcerns.vegan(), item.DietaryConcerns.vegetarian(), item.DietaryConcerns.paleo(), item.DietaryConcerns.glutenFree(), item.DietaryConcerns.kosher())

	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("failed to execute: %s", st)
	}
	item.Active = true
	err = put(ctx, id, item)
	if err != nil {
		return errDatastore.WithError(err).Wrapf("failed to put item(%d)", id)
	}
	return nil
}

func deleteActiveItem(ctx context.Context, id int64, item *Item) error {
	st := fmt.Sprintf(deleteStatement, id)
	_, err := mysqlDB.Exec(st)
	if err != nil {
		return errSQLDB.WithError(err).Wrapf("failed to execute: %s", st)
	}
	item.Active = false
	err = put(ctx, id, item)
	if err != nil {
		return errDatastore.WithError(err).Wrapf("failed to put item(%d)", id)
	}
	return nil
}

func connectSQL(ctx context.Context) {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = "root@/gigamunch"
	} else {
		projectID := config.GetProjectID(ctx)
		connectionString = fmt.Sprintf("root@cloudsql(%s:us-central1:gigasqldb)/gigamunch", projectID)
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}

type closer interface {
	Close() error
}

func handleCloser(ctx context.Context, c closer) {
	err := c.Close()
	if err != nil {
		utils.Errorf(ctx, "Error closing rows: %v", err)
	}
}
