package recipe

type Ingredient struct {
  Quantity string `json:"quantity"`
  Name     string `json:"name"`
  Prep     string `json:"prep,omitempty"`
  Note     string `json:"note,omitempty"`
}
