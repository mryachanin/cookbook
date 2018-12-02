package main

import (
  "gopkg.in/yaml.v2"
	"html/template"
  "io/ioutil"
	"log"
	"net/http"
	"path"
  "path/filepath"
  "./yaml/recipe"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	recipes, err := getRecipes()
  if (err != nil) {
    http.Error(w, err.Error(), 500)
    return
  }
	t := template.Must(template.ParseFiles(path.Join("templates", "index.html")))
	t.Execute(w, recipes)
}

func getRecipes() ([]recipe.Recipe, error) {
  recipes := []recipe.Recipe{}
  paths, err := filepath.Glob("public/recipes/*")
  if (err != nil) {
    return nil, err
  }
  for _, path := range paths {
    b, err := ioutil.ReadFile(path)
  	if err != nil {
      return nil, err
  	}

    var r recipe.Recipe
    if err = yaml.Unmarshal(b, &r); err != nil {
  		return nil, err
  	}
    recipes = append(recipes, r)
  }
  return recipes, nil
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
