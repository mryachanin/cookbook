package recipe

type Recipe struct {
  Name string
  Ingredients []Ingredient
  Instructions []string
  Yields string
  Keeps_for string
  Prep_time string
  Total_time string
  Notes []string
  Tags []string
  Source Source
}
