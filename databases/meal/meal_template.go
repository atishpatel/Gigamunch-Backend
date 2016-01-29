package meal

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/atishpatel/Gigamunch-Backend/databases/session"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// SaveMealTemplate saves the meal template if the user has permission
func SaveMealTemplate(ctx context.Context, sessionID string, mealTemplateID int64, mealTemplate *types.MealTemplate) (int64, *types.MealTemplate, error) {
	var err error
	var user *types.User
	var oldMealTemplate *types.MealTemplate
	// get user session
	userChannel := session.GetUserSession(ctx, sessionID)

	// not a new meal template
	if mealTemplateID != 0 {
		// get old meal template
		oldMealTemplate = &types.MealTemplate{}
		err = getMealTemplate(ctx, mealTemplateID, oldMealTemplate)
		user = <-userChannel
		if user == nil {
			err := errors.ErrSessionNotFound.WithArgs(sessionID)
			utils.Errorf(ctx, "%+v", err)
			return 0, nil, err
		}
		if err != nil {
			err := errors.ErrDatastore.WithArgs("get", "old meal template", user.Email, err)
			utils.Errorf(ctx, "%+v", err)
			return 0, nil, err
		}
		// check if user has right to access the meal template
		if user.Email != oldMealTemplate.GigachefEmail {
			err := errors.ErrUnauthorizedAccess.WithArgs(user.Email, "MealTemplate", mealTemplateID)
			utils.Errorf(ctx, "%+v", err)
			return 0, nil, err
		}
	} else {
		user = <-userChannel
		if user == nil {
			err := errors.ErrSessionNotFound.WithArgs(sessionID)
			utils.Errorf(ctx, "%+v", err)
			return 0, nil, err
		}
		// set up a default mealTemplate
		oldMealTemplate = &types.MealTemplate{
			Title: "Yummy Food",
			BaseMeal: types.BaseMeal{
				GigachefEmail: user.Email,
			},
			// TODO: Set default image and stuff
		}
	}
	oldOrginizationTags := oldMealTemplate.OrginizationTags
	newOrginizationTags := mealTemplate.OrginizationTags
	updatedMealTemplate := oldMealTemplate.ValidateAndUpdate(mealTemplate)
	mealTemplateID, err = saveMealTemplate(ctx, mealTemplateID, updatedMealTemplate)
	if err != nil {
		err := errors.ErrDatastore.WithArgs("put", "updated meal template", user.Email, err)
		utils.Errorf(ctx, "%+v", err)
		return 0, nil, err
	}
	updateOrginizationTags(ctx, oldOrginizationTags, newOrginizationTags, updatedMealTemplate)
	return mealTemplateID, updatedMealTemplate, nil
}

func getMealTemplate(ctx context.Context, mealTemplateID int64, mealTemplate *types.MealTemplate) error {
	// var err error
	if mealTemplateID == 0 {
		return datastore.ErrNoSuchEntity
	}
	mealTemplateKey := datastore.NewKey(ctx, types.KindMealTemplate, "", mealTemplateID, nil)
	return datastore.Get(ctx, mealTemplateKey, mealTemplate)
}

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
func GetMealTemplate(ctx context.Context, sessionID string, mealTemplateID int64) (*types.MealTemplate, error) {
	if mealTemplateID == 0 {
		err := errors.ErrInvalidParameter.WithArgs(0, "mealTemplateID")
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	mealTemplate := &types.MealTemplate{}
	userChannel := session.GetUserSession(ctx, sessionID)
	err := getMealTemplate(ctx, mealTemplateID, mealTemplate)
	user := <-userChannel
	if user == nil {
		err := errors.ErrSessionNotFound.WithArgs(sessionID)
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	if err != nil {
		err := errors.ErrDatastore.WithArgs("get", "old meal template", user.Email, err)
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	if user.Email != mealTemplate.GigachefEmail {
		err := errors.ErrUnauthorizedAccess.WithArgs(user.Email, "MealTemplate", mealTemplateID)
		utils.Errorf(ctx, "%+v", err)
		return nil, err
	}
	return mealTemplate, nil
}

// func ArchiveMealTemplate() {
//
// }
