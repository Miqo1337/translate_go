package main

import (
	"os"
	"path/filepath"
	"strings"
)

func processDir(dir string) error {
	return filepath.WalkDir(dir, func(path string, dirEnt os.DirEntry, err error) error {

		_, file := filepath.Split(path)

		if dirEnt.IsDir() {
			if file == "node_modules" || strings.Contains(file, "_dist") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(file) == ".js" {
			//fmt.Printf("test\n")
			return nil
		}

		processFile(path) //check to see if any difference between file as and argument or path as an argument

		return nil

	})
}
