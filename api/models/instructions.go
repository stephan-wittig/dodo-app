package models

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"text/template"
)

// InstructionSet is a request for or response with variable values
//
// This is created by DocumentTemplate and values are used to create a Document
type InstructionSet struct {
	TemplateId      string `json:"templateId"`
	TemplateVersion string `json:"templateVersion"`

	Elements  map[string]string              `json:"elements"`  // map of section/paragraph keys => label
	Variables map[string]variableInstruction `json:"variables"` // map of variable keys => variableInstructions
}

// variableInstruction is a generic set of instructions for one variable.
//
// This is supposed to be used in input prompts, not for evaluating a document
type variableInstruction struct {
	Value       string         `json:"value,omitempty"`    // Optional in requests. Everything else is optional in responses
	DataType    string         `json:"dataType,omitempty"` // "STRING", "CHOICE", or "NUMBER"
	Label       string         `json:"label,omitempty"`
	Description string         `json:"description,omitempty"`
	Min         float64        `json:"min,omitempty"`
	Max         float64        `json:"max,omitempty"`
	Step        float64        `json:"step,omitempty"`    // Only relevant for NUMBER
	Options     []choiceOption `json:"options,omitempty"` // Only relevant for CHOICE
}

// Constructors

// ParseJsonInstructionSet reads JSON data to construct an InstructionSet.
// NB: This does not validate the instructions
func ParseJsonInstructionSet(data []byte) (InstructionSet, error) {
	newInstructions := InstructionSet{}
	err := json.Unmarshal(data, &newInstructions)
	return newInstructions, err
}

// Methods

// Validate validates an InstructionSet using the relevant template
//
// Instructions are considered invalid if toggles for optional paragraphs or sections are missing,
// if variable values are missing, if variables don't comply with the template's constraints,
// or if StringVariables contain any XML elements
func (inst *InstructionSet) Validate() (bool, []error) {
	return true, []error{}
}

// GetDigest computes to MD5 digest of an InstructionSet
func (inst *InstructionSet) getDigest() ([16]byte, error) {
	jsonInstructions, err := json.Marshal(inst)
	if err != nil {
		return [16]byte{}, err
	}
	return md5.Sum(jsonInstructions), nil
}

// CreateDocument creates (or copies) a Document with a valid InstructionSet
func (inst *InstructionSet) CreateDocument() (Document, error) {
	intDoc, err := inst.createIntermediateDocument()
	if err != nil {
		return Document{}, err
	}

	digest, err := inst.getDigest()
	if err != nil {
		return Document{}, err
	}

	doc, err := intDoc.GenerateVerbatim()
	if err != nil {
		return Document{}, err
	}

	return Document{
		Id:             "",
		Digest:         digest,
		Valid:          false,
		Errors:         nil,
		Verbatim:       []byte(doc),
		InstructionSet: *inst,
	}, nil
}

// getTemplate retrieves the template referenced by this InstructionSet
func (inst *InstructionSet) getTemplate() (*DocumentTemplate, error) {
	return DemoTemplate, nil
}

func (inst *InstructionSet) createIntermediateDocument() (intermediateDocument, error) {
	tmpl, err := inst.getTemplate()
	if err != nil {
		return intermediateDocument{}, err
	}

	doc := intermediateDocument{
		Name:     tmpl.Name,
		Preambel: template.HTMLEscapeString(replaceGlobalVariables(tmpl.Preambel)),
		Sections: []intermediateSection{},
	}

	sectionCount := 0
	for _, section := range tmpl.Sections {
		if section.Optional && inst.Elements[section.Key] == "" {
			continue
		}

		sectionCount++

		sec := intermediateSection{
			// TODO: use roman numerals here
			Heading:    fmt.Sprintf("%d. %s", sectionCount, template.HTMLEscapeString(section.Heading)),
			Preambel:   template.HTMLEscapeString(replaceGlobalVariables(section.Preambel)),
			Paragraphs: []string{},
		}

		paragraphCount := 0

		for _, p := range section.Paragraphs {
			if p.Optional && inst.Elements[p.Key] == "" {
				continue
			}

			paragraphCount++

			copy, err := p.replaceAllVariables(*inst)
			copy = fmt.Sprintf("%d. %s", paragraphCount, copy)
			if err != nil {
				return intermediateDocument{}, err
			}
			sec.Paragraphs = append(sec.Paragraphs, copy)
		}

		doc.Sections = append(doc.Sections, sec)
	}

	return doc, nil
}
