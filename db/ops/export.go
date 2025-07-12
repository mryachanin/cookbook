// This contains logic to export all recipes from CouchDB to YAML files
// for backup purposes or data migration.
package ops

import (
  "github.com/mryachanin/cookbook/api/recipe"
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "github.com/rhinoman/couchdb-go"
  "gopkg.in/yaml.v3"
  "os"
  "log"
  "path/filepath"
  "strings"
  "regexp"
)

func ExportRecipes(outputPath string) {
  // Create output directory if it doesn't exist
  if err := os.MkdirAll(outputPath, 0755); err != nil {
    log.Fatalf("Failed to create output directory: %v", err)
  }

  c := config.LoadConfiguration(configPath)
  database := db.Connect(c)
  recipes := getAllRecipes(database)
  
  if len(recipes) == 0 {
    log.Printf("No recipes found in database")
    return
  }

  saveRecipesToFiles(outputPath, recipes)
  log.Printf("Successfully exported %d recipes to %s", len(recipes), outputPath)
}

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

func getAllRecipes(db *couchdb.Database) []recipe.Recipe {
  recipes := []recipe.Recipe{}

  // Get all recipe IDs using the existing view
  result := viewResult{}
  if err := db.GetView("recipe", "getRecipes", &result, nil); err != nil {
    log.Fatalf("Failed to get recipes from view: %v", err)
  }

  // Fetch each recipe by ID
  for _, row := range result.Rows {
    var r recipe.Recipe
    _, err := db.Read(row.Id, &r, nil)
    if err != nil {
      log.Printf("Failed to read recipe %s (%s): %v", row.Key, row.Id, err)
      continue
    }

    // Validate that it's actually a recipe (has a name)
    if r.Name == "" {
      log.Printf("Skipping document %s: missing recipe name", row.Id)
      continue
    }

    recipes = append(recipes, r)
  }

  return recipes
}

func saveRecipesToFiles(outputPath string, recipes []recipe.Recipe) {
  for _, r := range recipes {
    // Create a safe filename from the recipe name
    filename := sanitizeFilename(r.Name) + ".yaml"
    filePath := filepath.Join(outputPath, filename)

    // Clear internal CouchDB fields for clean export
    exportRecipe := r
    exportRecipe.Id = ""
    exportRecipe.Rev = ""

    // Marshal recipe to YAML
    yamlData, err := yaml.Marshal(&exportRecipe)
    if err != nil {
      log.Printf("Failed to marshal recipe '%s': %v", r.Name, err)
      continue
    }

    // Write to file
    if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
      log.Printf("Failed to write recipe '%s' to file: %v", r.Name, err)
      continue
    }

    log.Printf("Exported recipe: %s -> %s", r.Name, filename)
  }
}

func sanitizeFilename(name string) string {
  // Replace invalid filename characters with underscores
  reg := regexp.MustCompile(`[<>:"/\\|?*]`)
  sanitized := reg.ReplaceAllString(name, "_")
  
  // Replace spaces with underscores
  sanitized = strings.ReplaceAll(sanitized, " ", "_")
  
  // Remove multiple consecutive underscores
  reg = regexp.MustCompile(`_+`)
  sanitized = reg.ReplaceAllString(sanitized, "_")
  
  // Trim underscores from start and end
  sanitized = strings.Trim(sanitized, "_")
  
  // Ensure it's not empty
  if sanitized == "" {
    sanitized = "unnamed_recipe"
  }
  
  return sanitized
}