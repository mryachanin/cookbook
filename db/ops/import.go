// This contains logic to populate CouchDB with initial data for testing
// purposes until an API is implemented for creating recipes.
// This command is idempotent.
package ops

import (
  "github.com/mryachanin/cookbook/api/recipe"
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "github.com/google/uuid"
  "github.com/rhinoman/couchdb-go"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
  "path/filepath"
)

func ImportRecipes(path string) {
  recipesPath := path + "/*.yaml"

  recipes := getRecipes(recipesPath)
  if len(recipes) == 0 {
    log.Fatalf("Found no YAML recipe files at path: \"%s\"", recipesPath)
  }

  c := config.LoadConfiguration(configPath)
  db := db.Connect(c)
  storeRecipes(db, recipes)
}

func getRecipes(path string) ([]recipe.Recipe) {
  recipes := []recipe.Recipe{}

  // Get all of the recipe paths in the given directory.
  paths, err := filepath.Glob(path)
  if err != nil {
    log.Fatal(err)
  }

  // For each recipe in the given path, read and unmarshal it.
  for _, path := range paths {
    // Read the recipe.
    b, err := ioutil.ReadFile(path)
    if err != nil {
      log.Fatal(err)
    }

    // Unmarshal the recipe.
    var r recipe.Recipe
    if err = yaml.Unmarshal(b, &r); err != nil {
      log.Fatal(err)
    }

    // Append the recipe to the list of recipes.
    recipes = append(recipes, r)
  }
  return recipes
}

func storeRecipes(db *couchdb.Database, recipes []recipe.Recipe) {
  for _, v := range recipes {
    // Third argument is the revision. It is supposed to be blank on save.
    rev, err := db.Save(v, uuid.New().String(), "")
    if err != nil {
      // This command is idempotent. Don't fail if the recipe already exists.
      if err.(*couchdb.Error).StatusCode == 409 {
        log.Printf("Recipe with name \"%s\" already exists. Skipping.", v.Name)
      } else {
        log.Fatal(err)
      }
    } else {
      log.Printf("Saved recipe with name \"%s\". Rev: \"%s\"", v.Name, rev)
    }
  }
}
