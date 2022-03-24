package utils

import (
	"fmt"
	"io/fs"
	"os"
)

// OpenFile reads the file in the specified directory matching the given Glob.
// If multiple files match, the user is asked to select one
func OpenFile(dirName string, glob string) ([]byte, error) {
	if err := os.Chdir(dirName); err != nil {
		return nil, err
	}

	dir := os.DirFS(".")
	files, err := fs.Glob(dir, glob)
	if err != nil {
		return []byte{}, fmt.Errorf("Using glob failed: %s", err)
	}
	if files == nil {
		return []byte{}, fmt.Errorf("No file matching %s found", glob)
	}

	var file string
	numOfFiles := len(files)
	if numOfFiles == 1 {
		file = files[0]
	} else {
		var selection int
		fmt.Printf("Found %d files. Please choose one.\n", numOfFiles)
		for i, f := range files {
			fmt.Printf("  %2d: %s\n", i+1, f)
		}

		if _, err := fmt.Scan(&selection); err != nil {
			return []byte{}, fmt.Errorf("Could not read input: %s", err)
		}
		if selection > numOfFiles {
			return []byte{}, fmt.Errorf("Could not use input: Out of bounds")
		}
		file = files[selection-1]
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't read file: %s", err)
	}
	return data, nil
}
