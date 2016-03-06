package item

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/auth"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// SaveMealTemplate creates/updates the meal template if the user has permission
func SaveMealTemplate(ctx context.Context, authToken *types.AuthToken, mealTemplateID int64, mealTemplate *types.MealTemplate) (int64, error) {
	var err error
	var oldMealTemplate *types.MealTemplate
	// get user
	err = auth.GetUserFromToken(ctx, authToken)
	if err != nil {
		return 0, err
	}
	user := &authToken.User

	// not a new meal template
	if mealTemplateID != 0 {
		// get old meal template
		oldMealTemplate = &types.MealTemplate{}
		errChan := getMealTemplate(ctx, mealTemplateID, oldMealTemplate)
		err = <-errChan //TODO update to pass in errchan and check for ok
		if err != nil {
			err = errors.ErrDatastore.WithArgs("get", "old meal template", user.UserID, err)
			utils.Errorf(ctx, "%+v", err)
			return 0, err
		}
		// check if user has right to access the meal template
		if user.UserID != oldMealTemplate.GigachefID {
			err = errors.ErrUnauthorizedAccess.WithArgs(user.UserID, "MealTemplate", mealTemplateID)
			utils.Errorf(ctx, "%+v", err)
			return 0, err
		}
	} else {
		// set up a default mealTemplate
		oldMealTemplate = &types.MealTemplate{
			Title: "Yummy Food",
			BaseMeal: types.BaseMeal{
				GigachefID: user.UserID,
			},
			// TODO: Set default image and stuff
		}
	}
	oldOrginizationTags := oldMealTemplate.OrginizationTags
	newOrginizationTags := mealTemplate.OrginizationTags
	errs := mealTemplate.Validate()
	if errs.HasErrors() {
		return 0, errors.ErrInvalidParameter.WithArgs(errs.Error(), mealTemplate)
	}
	mealTemplateID, err = saveMealTemplate(ctx, mealTemplateID, mealTemplate)
	if err != nil {
		err = errors.ErrDatastore.WithArgs("put", "updated meal template", user.UserID, err)
		utils.Errorf(ctx, "%+v", err)
		return 0, err
	}
	updateOrginizationTags(ctx, oldOrginizationTags, newOrginizationTags, mealTemplate)
	return mealTemplateID, nil
}

// TODO add redis cache layer
// getMealTemplate is a helper function that manages getting from cache or database
func getMealTemplate(ctx context.Context, mealTemplateID int64, mealTemplate *types.MealTemplate) <-chan error {
	errChan := make(chan error)
	go func(ctx context.Context, mealTemplateID int64, mealTemplate *types.MealTemplate, errChan chan<- error) {
		defer close(errChan)
		if mealTemplateID == 0 {
			errChan <- datastore.ErrNoSuchEntity
			return
		}
		mealTemplateKey := datastore.NewKey(ctx, types.KindMealTemplate, "", mealTemplateID, nil)
		errChan <- datastore.Get(ctx, mealTemplateKey, mealTemplate)
	}(ctx, mealTemplateID, mealTemplate, errChan)
	return errChan
}

// TODO add redis cache layer
// saveMealTemplate is a helper function that manages saving in cache and database
func saveMealTemplate(ctx context.Context, mealTemplateID int64, mealTemplate *types.MealTemplate) (int64, error) {
	var mealTemplateKey *datastore.Key
	var err error
	if mealTemplateID == 0 {
		mealTemplateKey = datastore.NewIncompleteKey(ctx, types.KindMealTemplate, nil)
	} else {
		mealTemplateKey = datastore.NewKey(ctx, types.KindMealTemplate, "", mealTemplateID, nil)
	}
	mealTemplateKey, err = datastore.Put(ctx, mealTemplateKey, mealTemplate)
	return mealTemplateKey.IntID(), err
}

func updateOrginizationTags(ctx context.Context, oldOrginizationTags []string, newOrginizationTags []string, updatedMealTemplate *types.MealTemplate) error {
	//TODO(critical): remove and add MealTemplate info to tags
	return nil
}

// GetMealTemplate gets a meal template
// all errors from this function contain a HTTPStatusCode
func GetMealTemplate(ctx context.Context, authToken *types.AuthToken, mealTemplateID int64) (*types.MealTemplate, error) {
	var err error
	if mealTemplateID == 0 {
		err = errors.ErrInvalidParameter.WithArgs(0, "mealTemplateID")
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	err = auth.GetUserFromToken(ctx, authToken)
	if err != nil {
		return nil, err
	}
	user := &authToken.User
	mealTemplate := &types.MealTemplate{}
	errChan := getMealTemplate(ctx, mealTemplateID, mealTemplate)
	err = <-errChan // TODO update to add ok check and chan passed in
	if err != nil {
		err = errors.ErrDatastore.WithArgs("get", "old meal template", user.UserID, err)
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	if user.UserID != mealTemplate.GigachefID {
		err = errors.ErrUnauthorizedAccess.WithArgs(user.UserID, "MealTemplate", mealTemplateID)
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	return mealTemplate, nil
}

// func ArchiveMealTemplate() {
//
// }
