package meal

import (
	// mysql database
	"database/sql"
	"fmt"
	"log"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/docker/distribution/registry/api/errcode"
	// driver for mysql
	_ "github.com/go-sql-driver/mysql"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/types/queries"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

var (
	mysqlDB *sql.DB
)

// PostMeal posts a live meal if the post is valid
func PostMeal(ctx context.Context, sessionID string, meal *types.Meal) (int64, *types.Meal, errcode.Error) {
	var returnErr errcode.Error
	//var err error

	// get user session
	// userChannel := session.GetUserSession(ctx, sessionID)

	mealKey := datastore.NewIncompleteKey(ctx, types.KindMeal, nil)

	return mealKey.IntID(), meal, returnErr
}

func postMeal(ctx context.Context, mealTemplateID int64, meal *types.Meal) (int64, *types.Meal, error) {

	return 0, nil, nil
}

// getLiveMeals gets the live meal ids from the live meals database (mysql) and
// then calls getMultiMeal to get the actual meal information
func getLiveMeals(ctx context.Context, geopoint *types.GeoPoint, limitRange *types.LimitRange, sortedLiveMeals []*types.Meal, distances []float32) <-chan error {
	errChan := make(chan error)
	go func(ctx context.Context, geopoint *types.GeoPoint, limitRange *types.LimitRange, sortedLiveMeals []*types.Meal, distances []float32, errChan chan<- error) {
		defer close(errChan)
		var err error
		listLength := limitRange.EndLimit - limitRange.StartLimit
		if len(sortedLiveMeals) != listLength {
			errChan <- errors.ErrInvalidParameter.WithArgs("The size of list ", sortedLiveMeals)
			return
		}
		if len(distances) != listLength {
			errChan <- errors.ErrInvalidParameter.WithArgs("The size of list ", distances)
			return
		}
		liveMealQuery := queries.GetSortByDistanceQuery(geopoint.Latitude, geopoint.Longitude,
			types.DefaultMaxDistanceInMiles, limitRange.StartLimit, limitRange.EndLimit)
		rows, err := mysqlDB.Query(liveMealQuery)
		defer rows.Close()
		if err != nil {
			errChan <- err
			return // error
		}
		// TODO delete
		cols, _ := rows.Columns()
		utils.Debugf(ctx, "mysql query result columns: ", cols)
		mealIDs := make([]int64, listLength)
		for i := 0; rows.Next(); i++ {
			rows.Scan(&mealIDs[i], &distances[i])
		}
		err = rows.Err()
		if err != nil {
			errChan <- err
			return // error
		}
		err = getMultiMeal(ctx, mealIDs, sortedLiveMeals)
		errChan <- err
		return // success
	}(ctx, geopoint, limitRange, sortedLiveMeals, distances, errChan)
	return errChan
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
	return err
}

// GetLiveMeals gets live meals based on distance
func GetLiveMeals(ctx context.Context, sessionID string, geopoint *types.GeoPoint, limitRange *types.LimitRange) ([]*types.Meal, []float32, errcode.Error) {
	var returnErr errcode.Error
	var err error
	if !limitRange.Valid() {
		returnErr = errors.ErrInvalidParameter.WithArgs("Invalid range", limitRange)
		return nil, nil, returnErr
	}
	// make slices with correct sizes
	listLength := limitRange.EndLimit - limitRange.StartLimit
	liveMeals := make([]*types.Meal, listLength)
	sortedDistances := make([]float32, listLength)
	// get the actual live meals
	errChan := getLiveMeals(ctx, geopoint, limitRange, liveMeals, sortedDistances)
	err = <-errChan
	if err != nil {
		return nil, nil, errors.GetErrorWithStatusCode(err)
	}
	return liveMeals, sortedDistances, returnErr
}

// func GetFilteredLiveMeals(){}

func init() {
	var err error
	liveMealConfig := config.GetLiveMealConfig()

	var connectionString string
	if appengine.IsDevAppServer() {
		// "user:password@/dbname"
		connectionString = fmt.Sprintf("%s:%s@%s/gigamunch", liveMealConfig.MySQLUser,
			liveMealConfig.MySQLUserPassword, liveMealConfig.MySQLServerIP)
	} else {
		bgContext := context.Background()
		projectID := appengine.AppID(bgContext)
		instanceID := appengine.InstanceID()
		// "user@cloudsql(project-id:instance-name)/dbname"
		connectionString = fmt.Sprintf("%s@cloudsql(%s:%s)/gigamunch",
			liveMealConfig.MySQLUser, projectID, instanceID)
	}
	mysqlDB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Couldn't connect to mysql database")
	}
}
