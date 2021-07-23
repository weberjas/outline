package project

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/fatih/color"
)

// define regular expressions for elements of interest
var packageRegex = regexp.MustCompile(`(?m)^package (\S+)`)
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

	boldCyan := color.New(color.FgCyan).Add(color.Underline).Add(color.Bold)
	boldCyan.Printf("\n%s%s\n", calculateIndent(indentLevel), op.Name)
	if len(op.Packages) > 0 {
		for _, internalPackages := range op.Packages {
			internalPackages.Print(indentLevel + 1)
		}
	}

	for funcName := range op.Functions {
		color.Yellow("%s%s (F)\n", calculateIndent(indentLevel+1), op.Functions[funcName].Name)
	}

	for structName := range op.Structs {
		op.Structs[structName].Print(indentLevel + 1)
	}

}

func (os OutlineStruct) Print(indentLevel int) {
	color.Green("%s%s (S)\n", calculateIndent(indentLevel), os.Name)
	for methodName := range os.Methods {
		color.Magenta("%s%s (M)\n", calculateIndent(indentLevel+1), os.Methods[methodName].Name)
		// os.Methods[methodName].Print(indentLevel + 1)
	}
}

func calculateIndent(indentLevel int) string {
	// determine how far to indent
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent = indent + "    "
	}
	return indent
}

// parse the file contents and return the outlinePackage object created
func ParseFile(fileContents []byte, showOnlyExp bool) (OutlinePackage, error) {

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
			if unicode.IsUpper(rune(functionName[1][0])) || (unicode.IsLower(rune(functionName[1][0])) && !showOnlyExp) {
				projectContents.Functions[functionName[1]] = OutlineFunc{Name: functionName[1]}
			}
		}
	}
	// print structs
	if len(packageStructs) > 0 {
		for _, structName := range packageStructs {
			projectContents.Structs[structName[1]] = OutlineStruct{Name: structName[1], Methods: make(map[string]OutlineFunc)}

			var methodRegex = regexp.MustCompile(fmt.Sprintf(`\nfunc \(\S+ %s\) (\S+)[(]`, string(structName[1])))
			// extract struct methods
			packageMethods := methodRegex.FindAllStringSubmatch(string(fileContents), -1)

			for _, methodName := range packageMethods {
				if unicode.IsUpper(rune(methodName[1][0])) || (unicode.IsLower(rune(methodName[1][0])) && !showOnlyExp) {
					projectContents.Structs[structName[1]].Methods[methodName[1]] = OutlineFunc{Name: methodName[1]}
				}
			}
		}
	}
	return projectContents, nil

}
