package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var projectContents = make(map[string]outlinePackage)

// define regular expressions for elements of interest
var packageRegex = regexp.MustCompile(`^package (\S+)`)
var functionRegex = regexp.MustCompile(`func (\S+)[(]`)
var structRegex = regexp.MustCompile(`type (\S+) struct`)

type outlinePackage struct {
	name      string
	functions map[string]outlineFunc
	structs   map[string]outlineStruct
}

type outlineFunc struct {
	name    string
	calls   []string // TODO: this is going to be a hard one to figure out, do we just want to reference local functions?
	returns string
}

type outlineStruct struct {
	name  string
	funcs []outlineFunc
}

func main() {

	err := filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {

			fmt.Printf("\n-------------------\n")
			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Failed to read file: %s", path)
			}

			// extract the package name
			packageName := packageRegex.FindStringSubmatch(string(fileContents))[1]

			// extract package functions
			packageFunctions := functionRegex.FindAllStringSubmatch(string(fileContents), -1)

			// extract package structs
			packageStructs := structRegex.FindAllStringSubmatch(string(fileContents), -1)

			fmt.Printf("%s\n", packageName)
			// print functions
			if len(packageFunctions) > 0 {
				for _, functionName := range packageFunctions {
					fmt.Printf("F --> %s\n", string(functionName[1]))
				}
			}
			// print structs
			if len(packageStructs) > 0 {
				for _, structName := range packageStructs[1:] {
					fmt.Printf("S --> %s\n", string(structName[1]))
				}
			}

		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
