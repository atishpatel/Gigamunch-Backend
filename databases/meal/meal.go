package meal

import (
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

// PostMeal posts a live meal if the post is valid
func PostMeal(ctx context.Context, sessionID string, mealTemplateID int64, meal *types.Meal) (int64, *types.Meal, error) {

	var updatedMeal *types.Meal

	return mealTemplateID, updatedMeal, nil
}

// func GetMeal(){
//
// }

// GetLiveMeals gets live meals based on distance
func GetLiveMeals(ctx context.Context, geopoint *types.GeoPoint, startLimit int, endLimit int) ([]types.Meal, error) {
	var meals []types.Meal

	return meals, nil
}

// func GetFilteredLiveMeals(){}
