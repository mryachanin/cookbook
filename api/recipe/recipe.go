package recipe

type Recipe struct {
  Id           string   `json:"_id,omitempty"`
  Rev          string   `json:"_rev,omitempty"`
  Name         string
  Ingredients  []Ingredient
  Instructions []string
  Yields       string   `json:",omitempty"`
  Keeps_for    string   `json:",omitempty"`
  Prep_time    string   `json:",omitempty"`
  Total_time   string   `json:",omitempty"`
  Notes        []string `json:",omitempty"`
  Tags         []string `json:",omitempty"`
  Source       Source   `json:",omitempty"`
}
