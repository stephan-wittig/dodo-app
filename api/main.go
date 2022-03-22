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
	err := os.Chdir("../demo")
	if err != nil {
		log.Fatal("Changing WD failed:", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Opening WD failed:", err)
	}

	dir := os.DirFS("../demo")
	files, err := fs.Glob(dir, "*.2d2l")
	if err != nil {
		log.Fatal("Using glob failed:", err)
	}
	if files == nil {
		log.Fatal("No template definition found")
	}

	var file string
	numOfFiles := len(files)
	if numOfFiles == 1 {
		file = fmt.Sprintf("%s\\%s", wd, files[0])
	} else {
		var selection int
		fmt.Printf("Found %d template definitions. Please choose one.", numOfFiles)
		for i, f := range files {
			fmt.Printf("%2d: %s", i, f)
		}

		_, err := fmt.Scan(&selection)
		if err != nil {
			log.Fatal("Could not read input:", err)
		}
		if selection >= numOfFiles {
			log.Fatal("Could not read input:", err)
		}
		file = fmt.Sprintf("%s\\%s", wd, files[selection])
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("Couldn't read file:", err)
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
