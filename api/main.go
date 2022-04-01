package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/stephan-wittig/dodo/models"
	"github.com/stephan-wittig/dodo/utils"
)

func main() {
	templateData, err := utils.OpenFile("../demo", "*.dod")
	if err != nil {
		log.Fatal(err)
	}

	// Must be parsed to be stored in memory
	tmp, err := models.ParseDodDocumentTemplate(templateData)
	if err != nil {
		log.Fatal("Couldn't parse template:", err)
	}
	log.Println("Parsed template successfully!")

	fmt.Println("Do you want to (1) create instructions for the template or (2) generate a document from values?")

	var selection int64
	if _, err = fmt.Scan(&selection); err != nil {
		log.Fatalf("Could not read input: %s", err)
	}
	if selection > 2 || selection < 1 {
		log.Fatal("Could not read input: out of bounds")
	}

	if selection == 1 {
		instructions := tmp.CreateInstructionSet()

		instructionsJson, err := json.MarshalIndent(instructions, "", "  ")
		if err != nil {
			log.Fatal("Couldn't stringify instructions:", err)
		}

		// Save to file
		os.WriteFile(fmt.Sprintf("../demo/instructions_%s.json", tmp.Name), instructionsJson, os.ModePerm)
	}

	if selection == 2 {
		instructionsData, err := utils.OpenFile("../demo", "input_*.json")
		if err != nil {
			log.Fatal(err)
		}

		instructions, err := models.ParseJsonInstructionSet(instructionsData)
		if err != nil {
			log.Fatal("Couldn't parse instructions:", err)
		}
		log.Println("Parsed instructions successfully!")

		document, err := instructions.CreateDocument()
		if err != nil {
			log.Fatal("Couldn't generate document:", err)
		}
		log.Println("Created document sucessfully")

		fmt.Printf("Computed instructions hash:%x\n", document.Digest)
		// save to file
		os.WriteFile(fmt.Sprintf("../demo/%s.html", tmp.Name), document.Verbatim, os.ModePerm)
	}

	log.Println("Done")
}
