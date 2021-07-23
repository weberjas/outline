package project

import (
	"fmt"
	"regexp"
)

// define regular expressions for elements of interest
var packageRegex = regexp.MustCompile(`^package (\S+)`)
var functionRegex = regexp.MustCompile(`\nfunc (\S+)[(]`)
var structRegex = regexp.MustCompile(`\ntype (\S+) struct`)

type OutlinePackage struct {
	Name      string
	Packages  map[string]OutlinePackage
	Functions map[string]OutlineFunc
	Structs   map[string]OutlineStruct
}

type OutlineFunc struct {
	Name    string
	Calls   []string // TODO: this is going to be a hard one to figure out, do we just want to reference local functions?
	Returns string
}

type OutlineStruct struct {
	Name    string
	Methods map[string]OutlineFunc
}

func (op OutlinePackage) Print(indentLevel int) {
	fmt.Printf("Package: %s\n", op.Name)
	if len(op.Packages) > 0 {
		for _, internalPackages := range op.Packages {
			internalPackages.Print(indentLevel + 1)
		}
	}

	fmt.Printf("Functions:\n")

}
func (op OutlineFunc) Print() {

}
func (op OutlineStruct) Print() {

}

// parse the file contents and return the outlinePackage object created
func ParseFile(fileContents []byte) (OutlinePackage, error) {

	// extract the package name
	packageName := packageRegex.FindStringSubmatch(string(fileContents))[1]
	projectContents := OutlinePackage{Name: packageName, Functions: make(map[string]OutlineFunc), Structs: make(map[string]OutlineStruct)}

	// extract package functions
	packageFunctions := functionRegex.FindAllStringSubmatch(string(fileContents), -1)

	// extract package structs
	packageStructs := structRegex.FindAllStringSubmatch(string(fileContents), -1)

	// print functions
	if len(packageFunctions) > 0 {
		for _, functionName := range packageFunctions {
			projectContents.Functions[functionName[1]] = OutlineFunc{Name: functionName[1]}
			//fmt.Printf("F --> %s\n", string(functionName[1]))
		}
	}
	// print structs
	if len(packageStructs) > 0 {
		for _, structName := range packageStructs {
			//fmt.Printf("S --> %s\n", string(structName[1]))
			projectContents.Structs[structName[1]] = OutlineStruct{Name: structName[1], Methods: make(map[string]OutlineFunc)}

			var methodRegex = regexp.MustCompile(fmt.Sprintf(`\nfunc \(\S+ %s\) (\S+)[(]`, string(structName[1])))
			// extract struct methods
			packageMethods := methodRegex.FindAllStringSubmatch(string(fileContents), -1)

			for _, methodName := range packageMethods {
				//fmt.Printf("  M --> %s\n", string(methodName[1]))
				projectContents.Structs[structName[1]].Methods[methodName[1]] = OutlineFunc{Name: methodName[1]}
			}
		}
	}
	return projectContents, nil

}
