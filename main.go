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
var functionRegex = regexp.MustCompile(`\nfunc (\S+)[(]`)
var structRegex = regexp.MustCompile(`\ntype (\S+) struct`)

type outlinePackage struct {
	name      string
	packages  map[string]outlinePackage
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
	funcs map[string]outlineFunc
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
			projectContents[packageName] = outlinePackage{name: packageName, functions: make(map[string]outlineFunc), structs: make(map[string]outlineStruct)}

			// extract package functions
			packageFunctions := functionRegex.FindAllStringSubmatch(string(fileContents), -1)

			// extract package structs
			packageStructs := structRegex.FindAllStringSubmatch(string(fileContents), -1)

			fmt.Printf("%s\n", packageName)
			// print functions
			if len(packageFunctions) > 0 {
				for _, functionName := range packageFunctions {
					projectContents[packageName].functions[functionName[1]] = outlineFunc{name: functionName[1]}
					fmt.Printf("F --> %s\n", string(functionName[1]))
				}
			}
			// print structs
			if len(packageStructs) > 0 {
				for _, structName := range packageStructs {
					fmt.Printf("S --> %s\n", string(structName[1]))
					projectContents[packageName].structs[structName[1]] = outlineStruct{name: structName[1], funcs: make(map[string]outlineFunc)}

					var methodRegex = regexp.MustCompile(fmt.Sprintf(`\nfunc \(\S+ %s\) (\S+)[(]`, string(structName[1])))
					// extract struct methods
					packageMethods := methodRegex.FindAllStringSubmatch(string(fileContents), -1)

					for _, methodName := range packageMethods {
						fmt.Printf("  M --> %s\n", string(methodName[1]))
						projectContents[packageName].structs[structName[1]].funcs[methodName[1]] = outlineFunc{name: methodName[1]}

					}
				}
			}

		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
