package recipe

type Recipe struct {
  Id           string       `json:"_id,omitempty"`
  Rev          string       `json:"_rev,omitempty"`
  Name         string       `json:"name"`
  Ingredients  []Ingredient `json:"ingredients"`
  Instructions []string     `json:"instructions"`
  Yields       string       `json:"yields,omitempty"`
  KeepsFor     string       `json:"keeps_for,omitempty"  yaml:"keeps_for"`
  PrepTime     string       `json:"prep_time,omitempty"  yaml:"prep_time"`
  TotalTime    string       `json:"total_time,omitempty" yaml:"total_time"`
  Notes        []string     `json:"notes,omitempty"`
  Tags         []string     `json:"tags,omitempty"`
  Source       Source       `json:"source,omitempty"`
}
