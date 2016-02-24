package gigamuncher

// BaseMeal is the basic stuff in a Meal and MealTemplate
type BaseMeal struct {
	ID              int      `json:"id"`
	Title           string   `json:"title" endpoints:"req"`
	Description     string   `json:"description" endpoints:"req"`
	Photos          []string `json:"photos"`
	PricePerServing float32  `json:"price_per_serving" endpoints:"req"`
}

// Meal is a meal that is no longer live
type Meal struct {
	BaseMeal                // embedded
	ClosingDateTime int     `json:"closing_datetime" endpoints:"req"`
	ReadyDateTime   int     `json:"ready_datetime" endpoints:"req"`
	ServingsOffered int     `json:"servings_offered" endpoints:"req"`
	Distance        float32 `json:"distance"`
	OrdersLeft      int     `json:"orders_left"`
}

type MealDetailed struct {
	BaseMeal                  // embedded
	MealTemplateID   int      `json:"meal_template_id"`
	GeneralTags      []string `json:"general_tags"`
	CategoriesTags   []string `json:"categories_tags"`
	DietaryNeedsTags []string `json:"dietary_needs_tags"`
	NationalityTags  []string `json:"categories_tags"`
	Ingredients      []string `json:"ingredients"`
}
