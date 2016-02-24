package gigachef

import "time"

// BaseMeal is the basic stuff in a Meal and MealTemplate
type BaseMeal struct {
	Description              string   `json:"description" endpoints:"req"`
	CategoriesTags           []string `json:"categories_tags"`
	DietaryNeedsTags         []string `json:"dietary_needs_tags"`
	NationalityTags          []string `json:"categories_tags"`
	GeneralTags              []string `json:"general_tags"`
	Photos                   []string `json:"photos"`
	Ingredients              []string `json:"ingredients"`
	EstimatedPreperationTime int      `json:"estimated_preperation_time"`
	EstimatedCostPerServing  float32  `json:"estimated_cost_per_serving"`
	PricePerServing          float32  `json:"price_per_serving" endpoints:"req"`
}

// MealTemplate is the template Gigachefs can use to post a meal
type MealTemplate struct {
	BaseMeal                    // embedded
	Title             string    `json:"title"`
	LastUsedDataTime  time.Time `json:"lastused_datetime"`
	NumMealsCreated   int       `json:"num_meals_created"`
	NumTotalOrders    int       `json:"num_total_orders"`
	AverageMealRating float32   `json:"average_meal_rating"`
	OrginizationTags  []string  `json:"orginization_tags"`
}

// Meal is a meal that is no longer live
type Meal struct {
	BaseMeal               // embedded
	MealTemplateID  int    `json:"meal_template_id"`
	Title           string `json:"title" endpoints:"req"`
	ClosingDateTime int    `json:"closing_datetime" endpoints:"req"`
	ReadyDateTime   int    `json:"ready_datetime" endpoints:"req"`
	ServingsOffered int    `json:"servings_offered" endpoints:"req"`
}
