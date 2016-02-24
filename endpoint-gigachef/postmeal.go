package gigachef

import (
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/meal"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"golang.org/x/net/context"
)

type PostMealReq struct {
	Token string `json:"token" endpoints:"req"`
	Meal  Meal   `json:"meal" endpoints:"req"`
}

type PostMealResp struct {
	Token        string `json:"token"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	MealID       int    `json:"meal_id"`
}

// PostMeal is an endpoint that post a meal form a Gigachef
func (service *Service) PostMeal(ctx context.Context, req *PostMealReq) (*PostMealResp, error) {
	var err error
	authToken := &types.AuthToken{
		JWTString: req.Token,
	}
	newMeal := &types.Meal{
		BaseMeal: types.BaseMeal{
			Description:              req.Meal.BaseMeal.Description,
			CategoriesTags:           req.Meal.BaseMeal.CategoriesTags,
			DietaryNeedsTags:         req.Meal.BaseMeal.DietaryNeedsTags,
			NationalityTags:          req.Meal.BaseMeal.NationalityTags,
			GeneralTags:              req.Meal.BaseMeal.GeneralTags,
			Photos:                   req.Meal.BaseMeal.Photos,
			Ingredients:              req.Meal.BaseMeal.Ingredients,
			EstimatedPreperationTime: int64(req.Meal.BaseMeal.EstimatedPreperationTime),
			EstimatedCostPerServing:  req.Meal.BaseMeal.EstimatedCostPerServing,
			PricePerServing:          req.Meal.BaseMeal.PricePerServing,
		},
		MealTemplateID:  int64(req.Meal.MealTemplateID),
		Title:           req.Meal.Title,
		ClosingDateTime: time.Unix(int64(req.Meal.ClosingDateTime), 0),
		ReadyDateTime:   time.Unix(int64(req.Meal.ReadyDateTime), 0),
		ServingsOffered: req.Meal.ServingsOffered,
	}
	mealID, err := meal.PostMeal(ctx, authToken, newMeal)
	if err != nil {
		codeErr := errors.GetErrorWithStatusCode(err)
		resp := &PostMealResp{
			Token:        authToken.JWTString,
			ErrorCode:    codeErr.ErrorCode().Descriptor().HTTPStatusCode,
			ErrorMessage: codeErr.Error(),
		}
		return resp, err
	}
	resp := &PostMealResp{
		Token:  authToken.JWTString,
		MealID: int(mealID),
	}
	return resp, nil
}
