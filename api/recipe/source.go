package recipe

type Type string

const (
  Book Type    = "book"
  Website Type = "website"
)

type Source struct {
  Name string `json:"name"`
  Type Type   `json:"type"`
  Link string `json:"link"`
}
