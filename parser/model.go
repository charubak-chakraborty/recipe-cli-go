package parser

// Models Defined
type Output struct {
	UniqueRecipeCount       int                     `json:"unique_recipe_count"`
	CountPerRecipe          []CountPerRecipe        `json:"count_per_recipe"`
	BusiestPostcode         BusiestPostCode         `json:"busiest_postcode"`
	MatchByName             []string                `json:"match_by_name"`
	CountPerPostCodeAndTime CountPerPostCodeAndTime `json:"count_per_postcode_and_time"`
}
type CountPerRecipe struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

type BusiestPostCode struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type CountPerPostCodeAndTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}
