package view

import (
  appTemplate "github.com/mryachanin/cookbook/web/template"
  "github.com/rhinoman/couchdb-go"
  "html/template"
  "net/http"
)

type viewResult struct {
  TotalRows int              `json:"total_rows"`
  Offset    int              `json:"offset"`
  Rows      []viewResultItem `json:"rows"`
}

type viewResultItem struct {
  Id    string `json:"id"`
  Key   string `json:"key"`
  Value string `json:"value"`
}

func GetIndex(w http.ResponseWriter, r *http.Request, db *couchdb.Database) {
  result := viewResult{}
  if err := db.GetView("recipe", "getRecipes", &result, nil); err != nil {
    errorAndReturn(err, w)
  }

  templateArr := appTemplate.MustAsset("web/template/index.html")
  t, err := template.New("index.html").Parse(string(templateArr[:]))
  if err != nil {
    errorAndReturn(err, w)
  }
  if err := t.Execute(w, result.Rows); err != nil {
    errorAndReturn(err, w)
  }
}

func errorAndReturn(err error, w http.ResponseWriter) {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
}
