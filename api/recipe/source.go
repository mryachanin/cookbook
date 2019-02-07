package recipe

type Type string

const (
  Book Type = "book"
  Website Type = "website"
)

type Source struct {
  Name string
  Type Type
  Link string
}
