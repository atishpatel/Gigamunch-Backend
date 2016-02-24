package meal

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/docker/distribution/registry/api/errcode"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/core/account"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/types/queries"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

var (
	mysqlDB *sql.DB
)

// TODO add stuff to test for if verified chef

// PostMeal posts a live meal if the post is valid
// returns MealID, Meal, error
func PostMeal(ctx context.Context, authToken *types.AuthToken, meal *types.Meal) (int64, error) {
	var err error
	// get user
	err = auth.GetUserFromToken(ctx, authToken)
	if err != nil {
		return 0, err
	}
	meal.GigachefID = authToken.User.UserID
	// get the gigachef account for the location
	gigachefErrChan := make(chan error, 1)
	gigachef := &types.Gigachef{}
	account.GetGigachef(ctx, meal.GigachefID, gigachef, gigachefErrChan)
	err, ok := <-gigachefErrChan
	if !ok || err != nil {
		return 0, err
	}
	meal.Address = gigachef.Address
	// errs := meal.Validate()
	// if errs.HasErrors() {
	// 	return 0, errors.ErrInvalidParameter.WithArgs(errs.Error(), meal)
	// }
	return postMeal(ctx, meal)
}

func postMeal(ctx context.Context, meal *types.Meal) (int64, error) {
	var err error
	mealKey := datastore.NewIncompleteKey(ctx, types.KindMeal, nil)
	mealKey, err = datastore.Put(ctx, mealKey, meal)
	if err != nil {
		return 0, errors.ErrDatastore.WithArgs("put", "meal", meal.GigachefID, err)
	}
	mealID := mealKey.IntID()
	// TODO Put in mysql table
	_, err = mysqlDB.Exec(`INSERT
		INTO live_meals
		(meal_id, close_datetime, search_tags, is_experimental, is_baked_good, latitude, longitude)
		VALUES (?,?,?,?,?,?,?)`,
		mealID, meal.ClosingDateTime.Format(time.RFC3339), meal.Title, 0, 0, meal.Latitude, meal.Longitude)
	if err != nil {
		return 0, errors.ErrExternalDependencyFail.WithArgs("mysql", "inserting meal", err)
	}
	return mealID, nil
}

// getLiveMeals gets the live meal ids from the live meals database (mysql) and
// then calls getMultiMeal to get the actual meal information
func getLiveMeals(ctx context.Context, geopoint *types.GeoPoint, limitRange *types.LimitRange) ([]int64, []*types.Meal, []float32, error) {
	var err error
	listLength := limitRange.EndLimit - limitRange.StartLimit
	liveMealQuery := queries.GetSortByDistanceQuery(geopoint.Latitude, geopoint.Longitude,
		types.DefaultMaxDistanceInMiles, limitRange.StartLimit, limitRange.EndLimit)
	rows, err := mysqlDB.Query(liveMealQuery)
	if err != nil {
		return nil, nil, nil, err //TODO change to external dep err
	}
	defer rows.Close()
	tmpMealIDs := make([]int64, listLength)
	tmpDistances := make([]float32, listLength)
	actualReturnedRows := 0
	for i := 0; rows.Next(); i++ {
		rows.Scan(&tmpMealIDs[i], &tmpDistances[i])
		actualReturnedRows++
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, nil, err // TODO change to external dep err
	}
	// resize to actual meal returned size
	var mealIDs []int64
	var distances []float32
	if listLength != actualReturnedRows {
		mealIDs = make([]int64, actualReturnedRows)
		distances = make([]float32, actualReturnedRows)
		copy(mealIDs, tmpMealIDs)
		copy(distances, tmpDistances)
	} else {
		mealIDs = tmpMealIDs
		distances = tmpDistances
	}
	sortedLiveMeals := make([]*types.Meal, len(mealIDs))
	err = getMultiMeal(ctx, mealIDs, sortedLiveMeals)
	if err != nil {
		return nil, nil, nil, err
	}
	return mealIDs, sortedLiveMeals, distances, nil // success
}

// getMultiMeal gets a list of live meals from a list of meal IDs
// this function uses cache and datastore
func getMultiMeal(ctx context.Context, mealIDs []int64, dstMeals []*types.Meal) error {
	if mealIDs == nil || len(mealIDs) == 0 {
		return nil
	}
	if len(mealIDs) != len(dstMeals) {
		return fmt.Errorf("mealIDs and dstMeals slices have different length")
	}
	var err error
	// TODO add cache stuff
	mealKeys := make([]*datastore.Key, len(mealIDs))
	for i := range mealIDs {
		mealKeys[i] = datastore.NewKey(ctx, types.KindMeal, "", mealIDs[i], nil)
	}
	//TODO meals in dstMeals might be nil
	err = datastore.GetMulti(ctx, mealKeys, dstMeals)
	if err != nil && err.Error() != "(0 errors)" { // GetMulti always returns appengine.MultiError which is stupid
		return err
	}
	return nil
}

// GetLiveMeals gets live meals based on distance
func GetLiveMeals(ctx context.Context, geopoint *types.GeoPoint, limitRange *types.LimitRange) ([]int64, []*types.Meal, []float32, error) {
	var returnErr errcode.Error
	var err error
	if !limitRange.Valid() {
		returnErr = errors.ErrInvalidParameter.WithArgs("Invalid range", limitRange)
		return nil, nil, nil, returnErr
	}
	// get the actual live meals
	mealIDs, liveMeals, sortedDistances, err := getLiveMeals(ctx, geopoint, limitRange)
	utils.Debugf(ctx, "error in GetLiveMeals", err)
	if err != nil {
		return nil, nil, nil, errors.GetErrorWithStatusCode(err)
	}
	return mealIDs, liveMeals, sortedDistances, err
}

// func GetFilteredLiveMeals(){}

func init() {
	var err error
	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = fmt.Sprintf("root@/gigamunch")
	} else {
		// bgContext := context.Background()
		// projectID := appengine.AppID(bgContext)
		// "user@cloudsql(project-id:instance-name)/dbname"
		connectionString = fmt.Sprintf("root@cloudsql(gigamunch-omninexus-dev:gigasql)/gigamunch")
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}
