package models

import (
	"encoding/xml"

	"github.com/stephan-wittig/dodo/utils"
)

// Global variable holding the template used for everything
//
// TODO: Remove this and retrieve templates instead
var DemoTemplate *DocumentTemplate

// optionalElement composes the common attributes for optional document part templates
type optionalElement struct {
	Optional    bool   `xml:"optional,attr"` // If Optional = true, Description and Key must be set
	Description string `xml:"description,attr"`
	Key         string `xml:"key,attr"`
}

// DocumentTemplate is an internal representation of templates written in DOD
type DocumentTemplate struct {
	Id      string `xml:"id,attr"`
	Version string `xml:"version,attr"` // semver version (Major.Minor.Fix)

	Name         string `xml:"name,attr"`
	Description  string `xml:"description,attr"`
	Jurisdiction string `xml:"jurisdiction,attr"`
	Preambel     string `xml:"Preambel"` // Optional

	Sections []sectionTemplate `xml:"Section"`
}

type sectionTemplate struct {
	Heading string `xml:"heading,attr"`

	Subsections []subsectionTemplate `xml:"Subsection"`

	optionalElement
}

type subsectionTemplate struct {
	Label      string              `xml:"label,attr"` // Mandatory for optional Subsections
	Preambel   string              `xml:"Preambel"`   // Optional
	Paragraphs []paragraphTemplate `xml:"Paragraph"`

	optionalElement
}

type paragraphTemplate struct {
	Label    string `xml:"label,attr"` // Mandatory for optional paragraphs
	Verbatim string `xml:"Verbatim"`

	StringVariables []stringVariable `xml:"StringVariable"`
	NumberVariables []numberVariable `xml:"NumberVariable"`

	optionalElement
}

// Constructors

// ParseDocumentTemplate reads DOD data to construct a DocumentTemplate
//
// The DOD XML Schema is used to validate the input.
// On error, this returns the error AND a DocumentTemplate which might be invalid
func ParseDodDocumentTemplate(data []byte) (*DocumentTemplate, error) {
	newTemplate := DocumentTemplate{}
	trimmedData := utils.TrimNewLines(data)
	err := xml.Unmarshal(trimmedData, &newTemplate)
	DemoTemplate = &newTemplate
	return &newTemplate, err
}

// Methods

// CreateInstructionSet creates an InstructionSet for a given DocumentTemplate
func (tmpl *DocumentTemplate) CreateInstructionSet() InstructionSet {
	variables := map[string]variableInstruction{}
	elements := map[string]string{}

	for _, sec := range tmpl.Sections {
		if sec.Optional {
			elements[sec.Key] = sec.Heading
		}

		for _, sub := range sec.Subsections {
			if sub.Optional {
				elements[sub.Key] = sub.Label
			}

			for _, par := range sub.Paragraphs {
				if par.Optional {
					elements[par.Key] = par.Label
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

// ReplaceAllVariables replaces global and local variables
func (p *paragraphTemplate) replaceAllVariables(i InstructionSet) (string, error) {
	copy := p.Verbatim

	// First, replace locals
	for _, v := range p.StringVariables {
		ins, _ := i.Variables[v.Key]
		newCopy, err := replaceVariable(copy, v, ins.Value)
		if err != nil {
			return "", err
		}
		copy = newCopy
	}

	for _, v := range p.NumberVariables {
		ins, _ := i.Variables[v.Key]
		newCopy, err := replaceVariable(copy, v, ins.Value)
		if err != nil {
			return "", err
		}
		copy = newCopy
	}

	// Then, replace globals
	copy = replaceGlobalVariables(copy)

	return copy, nil
}
