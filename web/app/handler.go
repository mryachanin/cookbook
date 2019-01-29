package app

import (
	"gitlab.com/mryachanin/satisfied-vegan/api/recipe"
	appTemplate "gitlab.com/mryachanin/satisfied-vegan/web/template"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func HandleRequests() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	recipes, err := getRecipes()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	templateArr := appTemplate.MustAsset("web/template/index.html")
	t := template.New("index.html")
	t.Parse(string(templateArr[:]))
	t.Execute(w, recipes)
}

func getRecipes() ([]recipe.Recipe, error) {
	recipes := []recipe.Recipe{}
	paths, err := filepath.Glob("init/data/*")
	if err != nil {
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
