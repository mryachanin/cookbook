package view

import (
  appTemplate "github.com/mryachanin/cookbook/web/template"
  "github.com/julienschmidt/httprouter"
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

func GetIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  result := viewResult{}
  if err := db.GetView("recipe", "getRecipes", &result, nil); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  templateArr := appTemplate.MustAsset("web/template/index.html")
  t, err := template.New("index.html").Parse(string(templateArr[:]))
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if err := t.Execute(w, result.Rows); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}
