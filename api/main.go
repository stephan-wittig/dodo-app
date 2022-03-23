package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"stephan-wittig/dodo/models"
)

func main() {
	data, err := open2d2lFile()
	if err != nil {
		log.Fatal(err)
	}

	template, err := models.Parse2d2lDocumentTemplate(data)
	if err != nil {
		log.Fatal("Couldn't parse template:", err)
	}

	instructions := template.CreateInstructionSet()

	instructionsJson, err := json.Marshal(instructions)
	if err != nil {
		log.Fatal("Couldn't stringify instructions:", err)
	}

	fmt.Printf("Parsed template successfully!\n%s\n", instructionsJson)

	log.Println("Done")
}

func open2d2lFile() ([]byte, error) {
	err := os.Chdir("../demo")
	if err != nil {
		return []byte{}, fmt.Errorf("Changing WD failed: %s", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return []byte{}, fmt.Errorf("Opening WD failed: %s", err)
	}

	dir := os.DirFS("../demo")
	files, err := fs.Glob(dir, "*.2d2l")
	if err != nil {
		return []byte{}, fmt.Errorf("Using glob failed: %s", err)
	}
	if files == nil {
		return []byte{}, fmt.Errorf("No template definition found in ../demo")
	}

	var file string
	numOfFiles := len(files)
	if numOfFiles == 1 {
		file = fmt.Sprintf("%s\\%s", wd, files[0])
	} else {
		var selection int
		fmt.Printf("Found %d template definitions. Please choose one.\n", numOfFiles)
		for i, f := range files {
			fmt.Printf("  %2d: %s\n", i, f)
		}

		_, err := fmt.Scan(&selection)
		if err != nil {
			return []byte{}, fmt.Errorf("Could not read input: %s", err)
		}
		if selection >= numOfFiles {
			return []byte{}, fmt.Errorf("Could not use input: Out of bounds")
		}
		file = fmt.Sprintf("%s\\%s", wd, files[selection])
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return []byte{}, fmt.Errorf("Couldn't read file: %s", err)
	}
	return data, nil
}
