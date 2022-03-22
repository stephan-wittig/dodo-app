package models

import (
	"encoding/xml"
	"stephan-wittig/dodo/utils"
)

// optionalElement composes the common attributes for optional document part templates
type optionalElement struct {
	Optional    bool   `xml:"optional,attr"` // If Optional = true, Description and Key must be set
	Description string `xml:"description,attr"`
	Key         string `xml:"key,attr"`
}

// DocumentTemplate is an internal representation of templates written in 2D2L
type DocumentTemplate struct {
	Id      string `xml:"id,attr"`
	Version string `xml:"version,attr"` // semver version (Major.Minor.Fix)

	Name         string `xml:"name,attr"`
	Description  string `xml:"description,attr"`
	Jurisdiction string `xml:"jurisdiction,attr"`
	Preambel     string `xml:"Preambel"` // Optional

	Sections []sectionTemplate `xml:"SectionTemplate"`
}

type sectionTemplate struct {
	Heading  string `xml:"heading,attr"`
	Preambel string `xml:"Preambel"` // Optional

	Paragraphs []paragraphTemplate `xml:"ParagraphTemplate"`

	optionalElement
}

type paragraphTemplate struct {
	Label string `xml:"label,attr"` // Mandatory for optional paragraphs

	StringVariables []stringVariable `xml:"StringVariable"`
	NumberVariables []numberVariable `xml:"NumberVariable"`
	ChoiceVariables []choiceVariable `xml:"ChoiceVariable"`

	optionalElement
}

// Constructors

// ParseDocumentTemplate reads 2D2L data to construct a DocumentTemplate
//
// The 2D2L XML Schema is used to validate the input.
// On error, this returns the error AND a DocumentTemplate which might be invalid
func Parse2d2lDocumentTemplate(data []byte) (DocumentTemplate, error) {
	newTemplate := DocumentTemplate{}
	trimmedData := utils.TrimNewLines(data)
	err := xml.Unmarshal(trimmedData, &newTemplate)
	return newTemplate, err
}

// Methods

// CreateInstructionSet creates an InstructionSet for a given DocumentTemplate
func (tmpl *DocumentTemplate) CreateInstructionSet() InstructionSet {
	variables := map[string]variableInstruction{}
	elements := map[string]bool{}

	for _, sec := range tmpl.Sections {
		if sec.Optional {
			elements[sec.Key] = true
		}

		for _, par := range sec.Paragraphs {
			if par.Optional {
				elements[par.Key] = true
			}

			if par.StringVariables != nil {
				for _, v := range par.StringVariables {
					variables[v.GetKey()] = v.ToInstructions()
				}
			}

			if par.NumberVariables != nil {
				for _, v := range par.NumberVariables {
					variables[v.GetKey()] = v.ToInstructions()
				}
			}

			if par.ChoiceVariables != nil {
				for _, v := range par.ChoiceVariables {
					variables[v.GetKey()] = v.ToInstructions()
				}
			}
		}
	}

	return InstructionSet{
		TemplateId:      tmpl.Id,
		TemplateVersion: tmpl.Version,
		Elements:        elements,
		Variables:       variables,
	}
}
