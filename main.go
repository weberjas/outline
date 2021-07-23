package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/outline/project"
)

var projectContents = make(map[string]project.OutlinePackage)

func main() {

	showOnlyExp := false
	err := filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {

			fmt.Printf("Loading %s\n", path)
			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Failed to read file: %s", path)
			}
			parsedOutlinePackage, err := project.ParseFile(fileContents, showOnlyExp)
			if err != nil {
				log.Printf("Failed to parse file: %s", path)
			}
			projectContents[parsedOutlinePackage.Name] = parsedOutlinePackage
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// print out project
	for pkg := range projectContents {
		projectContents[pkg].Print(0)
	}
}
