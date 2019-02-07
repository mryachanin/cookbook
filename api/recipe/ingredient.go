package recipe

type Ingredient struct {
  Quantity string
  Name     string
  Prep     string `json:",omitempty"`
  Note     string `json:",omitempty"`
}
