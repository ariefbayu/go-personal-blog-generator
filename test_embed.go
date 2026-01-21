package main

import (
"embed"
"fmt"
"io/fs"
"log"
)

//go:embed admin-files/**
var adminFS embed.FS

func main() {
	// Create sub-filesystem
	subFS, err := fs.Sub(adminFS, "admin-files")
	if err != nil {
		log.Fatal(err)
	}
	
	// Try to read dashboard.html
	content, err := fs.ReadFile(subFS, "content/dashboard.html")
	if err != nil {
		log.Fatal("Error reading dashboard.html: ", err)
	}
	
	fmt.Printf("Successfully read %d bytes\n", len(content))
	
	// List all files
	fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
if err != nil {
return err
}
fmt.Println(path)
return nil
})
}
