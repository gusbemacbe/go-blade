package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gusbemacbe/go-blade"
)

func main() {
	// Creating a temporary directory for the compiled cache.
	tempCacheDir, err := os.MkdirTemp("", "blade-test-cache")

	if err != nil {
		log.Fatalf("Failed to create temp cache dir: %v", err)
	}
	defer os.RemoveAll(tempCacheDir)

	// Initializing our modernized Blade engine.
	// It looks for templates in the "views" directory.
	engine := blade.New([]string{"views"}, tempCacheDir)

	// Defining the data to be passed to the template.
	data := map[string]interface{}{
		"Name":    "Ben-Zod",
		"IsAdmin": true,
	}

	// Running the template engine to render the "test" view.
	html, err := engine.Run("test", data)

	if err != nil {
		log.Fatalf("Failed to render Blade template: %v", err)
	}

	// Printing the final rendered HTML to the console.
	fmt.Println("--- Rendered HTML Output ---")
	fmt.Println(html)
}
