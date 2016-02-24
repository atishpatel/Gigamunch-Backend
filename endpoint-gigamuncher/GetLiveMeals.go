package gigamuncher

import (
	"fmt"

	"github.com/atishpatel/Gigamunch-Backend/core/meal"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
	"golang.org/x/net/context"
)

// GetLiveMealsReq is the input required to get a list of live meals
type GetLiveMealsReq struct {
	StartLimit int     `json:"start_limit"`
	EndLimit   int     `json:"end_limit" endpoints:"req"`
	Latitude   float32 `json:"latitude" endpoints:"req"`
	Longitude  float32 `json:"longitude" endpoints:"req"`
	UserID     string  `json:"user_id"`
}

// Valid returns an error if input in invalid
func (req *GetLiveMealsReq) Valid() error {
	limit := types.LimitRange{StartLimit: req.StartLimit, EndLimit: req.EndLimit}
	if !limit.Valid() {
		return fmt.Errorf("Limit is not valid")
	}
	point := types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	if !point.Valid() {
		return fmt.Errorf("Location inputed is not valid")
	}
	return nil
}

// GetLiveMealsResp returns a list of meals
// TODO make a meal type that is returned
type GetLiveMealsResp struct {
	Meals []Meal `json:"meals"`
}

// GetLiveMeals is an endpoint that returns a list of meals
func (service *Service) GetLiveMeals(ctx context.Context, req *GetLiveMealsReq) (*GetLiveMealsResp, error) {
	var err error
	err = req.Valid()
	if err != nil {
		utils.Errorf(ctx, "There was an error with the input of: %+v", err)
		return nil, err
	}
	point := &types.GeoPoint{Latitude: req.Latitude, Longitude: req.Longitude}
	limit := &types.LimitRange{StartLimit: req.StartLimit, EndLimit: req.EndLimit}
	mealIDs, mealPointers, distances, err := meal.GetLiveMeals(ctx, point, limit)
	if err != nil {
		utils.Errorf(ctx, "There was an error getting live meals: %+v", err)
		return nil, err
	}
	meals := make([]Meal, len(mealPointers))
	for i := range mealPointers {
		meals[i].ID = int(mealIDs[i])
		meals[i].Distance = distances[i]
		meals[i].Title = mealPointers[i].Title
		meals[i].ClosingDateTime = int(mealPointers[i].ClosingDateTime.Unix())
		meals[i].ReadyDateTime = int(mealPointers[i].ReadyDateTime.Unix())
		meals[i].ServingsOffered = mealPointers[i].ServingsOffered
		meals[i].Description = mealPointers[i].Description
		meals[i].Photos = mealPointers[i].Photos
		meals[i].PricePerServing = mealPointers[i].PricePerServing
		meals[i].OrdersLeft = mealPointers[i].ServingsOffered - len(mealPointers[i].Orders)
	}
	resp := &GetLiveMealsResp{
		Meals: meals,
	}
	return resp, nil
}
