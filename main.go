package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/outline/project"
)

var projectContents = make(map[string][]project.OutlinePackage)

func main() {

	err := filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {

			//fmt.Printf("Loading %s\n", path)
			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Failed to read file: %s", path)
			}
			parsedOutlinePackage, err := project.ParseFile(fileContents)
			if err != nil {
				log.Printf("Failed to parse file: %s", path)
			}
			dir := filepath.Dir(path)
			projectContents[dir] = append(projectContents[dir], parsedOutlinePackage)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// print out project
	boldWhite := color.New(color.FgWhite).Add(color.Underline).Add(color.Bold)

	// print main if it exists
	mainPkg := projectContents["main"]
	if mainPkg != nil {
		mainPkg[0].Print(0)
		delete(projectContents, "main")
	}

	// sort the keys
	keys := make([]string, 0, len(projectContents))
	for k := range projectContents {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		pathParts := strings.Split(key, "/")
		indentCount := len(pathParts)
		indent := project.CalculateIndent(indentCount)

		boldWhite.Printf("\n\n%s%s\n", indent, pathParts[len(pathParts)-1])
		for _, file := range projectContents[key] {
			file.Print(indentCount + 1)
		}
	}
}
