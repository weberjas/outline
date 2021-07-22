package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func main() {
	err := filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		fmt.Println(info.Name())
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
