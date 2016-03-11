package gigachef

// BaseItem is the basic stuff in an Item
type BaseItem struct {
	Title            string   `json:"title"`
	Subtitle         string   `json:"subtitle"`
	Description      string   `json:"description"`
	Ingredients      []string `json:"ingredients"`
	GeneralTags      []string `json:"general_tags"`
	CuisineTags      []string `json:"cuisine_tags"`
	DietaryNeedsTags []string `json:"dietary_needs_tags"`
	Photos           []string `json:"photos"`
}
