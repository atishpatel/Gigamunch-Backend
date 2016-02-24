package types

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

// BaseMeal is the basic stuff in a Meal and MealTemplate
type BaseMeal struct {
	GigachefID               string    `json:"gigachef_id" datastore:",index"`
	CreatedDateTime          time.Time `json:"created_datetime" datastore:",noindex"`
	Description              string    `json:"description" datastore:",noindex"`
	GeneralTags              []string  `json:"general_tags" datastore:",noindex"`
	CategoriesTags           []string  `json:"categories_tags" datastore:",noindex"`
	DietaryNeedsTags         []string  `json:"dietary_needs_tags" datastore:",noindex"`
	NationalityTags          []string  `json:"categories_tags" datastore:",noindex"`
	Photos                   []string  `json:"photos" datastore:",noindex"`
	Ingredients              []string  `json:"ingredients" datastore:",noindex"`
	EstimatedPreperationTime int64     `json:"estimated_preperation_time" datastore:",noindex"`
	EstimatedCostPerServing  float32   `json:"estimated_cost_per_serving" datastore:",noindex"`
	PricePerServing          float32   `json:"price_per_serving" datastore:",noindex"`
}

// Validate validates the BaseMeal properties.
// The form is valid if errors.Errors.HasErrors() == false.
func (baseMeal *BaseMeal) Validate() errors.Errors {
	var multipleErrors errors.Errors
	if baseMeal.GigachefID == "" {
		multipleErrors.AddError(fmt.Errorf("GigachefID is empty"))
	}
	if baseMeal.CreatedDateTime.Year() < 3 {
		multipleErrors.AddError(fmt.Errorf("CreatedDateTime is not set"))
	}
	if len(baseMeal.Description) > 10 || utils.ContainsBanWord(baseMeal.Description) {
		multipleErrors.AddError(fmt.Errorf("Description is too short"))
	}
	if len(baseMeal.Photos) == 0 {
		multipleErrors.AddError(fmt.Errorf("Photos must be more than zero"))
	}
	if baseMeal.PricePerServing < 1 {
		multipleErrors.AddError(fmt.Errorf("PricePerServing must be more than $1.00"))
	}
	// TODO(Atish): add stuff like perp time if it's required
	return multipleErrors
}

// ValidateAndUpdate updates the BaseMeal based on the new meal template for the fields that pass validation
func (oldBaseMeal *BaseMeal) ValidateAndUpdate(newBaseMeal *BaseMeal) *BaseMeal {
	if oldBaseMeal.GigachefID == "" {
		oldBaseMeal.GigachefID = newBaseMeal.GigachefID
	}
	if oldBaseMeal.CreatedDateTime.Year() < 3 {
		oldBaseMeal.CreatedDateTime = time.Now().UTC()
	}
	if len(newBaseMeal.Description) > 10 && !utils.ContainsBanWord(newBaseMeal.Description) {
		oldBaseMeal.Description = newBaseMeal.Description
	}
	oldBaseMeal.CategoriesTags = newBaseMeal.CategoriesTags
	oldBaseMeal.DietaryNeedsTags = newBaseMeal.DietaryNeedsTags
	oldBaseMeal.NationalityTags = newBaseMeal.NationalityTags
	oldBaseMeal.GeneralTags = newBaseMeal.GeneralTags
	oldBaseMeal.Photos = newBaseMeal.Photos
	oldBaseMeal.Ingredients = newBaseMeal.Ingredients
	if newBaseMeal.EstimatedPreperationTime > 0 {
		oldBaseMeal.EstimatedPreperationTime = newBaseMeal.EstimatedPreperationTime
	}
	if newBaseMeal.EstimatedCostPerServing > -1 {
		oldBaseMeal.EstimatedCostPerServing = newBaseMeal.EstimatedCostPerServing
	}
	if newBaseMeal.PricePerServing > 0 {
		oldBaseMeal.PricePerServing = newBaseMeal.PricePerServing
	}
	return oldBaseMeal
}

// MealTemplate is the template Gigachefs can use to post a meal
type MealTemplate struct {
	BaseMeal                    // embedded
	Title             string    `json:"title" datastore:",index"`
	LastUsedDataTime  time.Time `json:"lastused_datetime" datastore:",index"`
	NumMealsCreated   int64     `json:"num_meals_created" datastore:",noindex"`
	NumTotalOrders    int64     `json:"num_total_orders" datastore:",noindex"`
	AverageMealRating float32   `json:"average_meal_rating" datastore:",index"`
	OrginizationTags  []string  `json:"orginization_tags" datastore:",noindex"`
}

// ValidateAndUpdate updates the MealTemplate based on the new MealTemplate for the fields that pass validation
func (oldMealTemplate *MealTemplate) ValidateAndUpdate(newMealTemplate *MealTemplate) *MealTemplate {
	if len(newMealTemplate.Title) > 3 && !utils.ContainsBanWord(newMealTemplate.Title) {
		oldMealTemplate.Title = newMealTemplate.Title
	}
	oldMealTemplate.LastUsedDataTime = time.Now().UTC()
	if newMealTemplate.NumMealsCreated > 0 {
		oldMealTemplate.NumMealsCreated = newMealTemplate.NumMealsCreated
	}
	if newMealTemplate.NumTotalOrders > 0 {
		oldMealTemplate.NumTotalOrders = newMealTemplate.NumTotalOrders
	}
	oldMealTemplate.AverageMealRating = newMealTemplate.AverageMealRating
	oldMealTemplate.OrginizationTags = newMealTemplate.OrginizationTags
	oldMealTemplate.BaseMeal = *oldMealTemplate.BaseMeal.ValidateAndUpdate(&newMealTemplate.BaseMeal)
	return newMealTemplate
}

// MealTemplateReference is a reference stored to a MealTemplate with basic rendering details
type MealTemplateReference struct {
	MealTemplateID   int64     `json:"meal_template_id" datstore:",noindex"`
	Title            string    `json:"title" datastore:",noindex"`
	PhotoURL         string    `json:"photo_url" datastore:",noindex"`
	Description      string    `json:"description" datastore:",noindex"`
	LastUsedDateTime time.Time `json:"lastused_datatime" datastore:",noindex"`
}

// MealTemplateTag is an orginizing structure for Gigachefs
type MealTemplateTag struct {
	Tag                    string                  `json:"tag" datastore:",index"`
	GigachefID             string                  `json:"gigachef_id" datastore:",index"`
	NumMealTemplates       int64                   `json:"num_meal_templates" datastore:",noindex"`
	MealTemplateReferences []MealTemplateReference `json:"meal_template_references" datastore:",noindex"`
}

// Meal is a meal that is no longer live
type Meal struct {
	BaseMeal                  // embedded
	MealTemplateID  int64     `json:"meal_template_id" datastore:",noindex"`
	Title           string    `json:"title" datastore:",noindex"`
	IsLive          bool      `json:"islive" datastore:",index"`
	ClosingDateTime time.Time `json:"closing_datetime" datastore:",noindex"`
	ReadyDateTime   time.Time `json:"ready_datetime" datastore:",index"`
	ServingsOffered int       `json:"servings_offered" datastore:",noindex"`
	TotalTips       float32   `json:"total_tips" datastore:",noindex"`
	TotalRevenue    float32   `json:"total_revenue" datastore:",noindex"`
	NumOrders       int64     `json:"num_orders" datastoer:",noindex"`
	Orders          []int64   `json:"orders" datastore:",noindex"`
	Address                   // embedded
}

// Validate validates the BaseMeal properties.
// The form is valid if errors.Errors.HasErrors() == false.
func (meal *Meal) Validate() errors.Errors {
	var multipleErrors errors.Errors
	if meal.ClosingDateTime.After(meal.ReadyDateTime) {
		multipleErrors.AddError(fmt.Errorf("ClosingDateTime cannot be after ReadyDateTime"))
	}
	if meal.ServingsOffered < 1 {
		multipleErrors.AddError(fmt.Errorf("ServingsOffered need to be greater than 0"))
	}
	multipleErrors = append(multipleErrors, meal.BaseMeal.Validate())
	return multipleErrors
}
