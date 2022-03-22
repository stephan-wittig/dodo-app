package models

import "encoding/json"

// InstructionSet is a request for or response with variable values
//
// This is created by DocumentTemplate and used to create a Document
type InstructionSet struct {
	TemplateId      string `json:"templateId"`
	TemplateVersion string `json:"templateVersion"`

	Elements  map[string]bool                `json:"elements"`  // map of section/paragraph keys => bool
	Variables map[string]variableInstruction `json:"variables"` // map of variable keys => variableInstructions
}

type variableInstruction struct {
	Value       string         `json:"value,omitempty"`       // Optional in requests. Everything else is optional in responses
	DataType    string         `json:"dataType,omitempty"`    // "STRING", "CHOICE", or "NUMBER"
	Description string         `json:"description,omitempty"` // Optional
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
	return false, []error{}
}

// GetDigest computes to MD5 digest of an InstructionSet
func (inst *InstructionSet) getDigest() ([]byte, error) {
	return []byte{}, nil
}

// CreateDocument creates (or copies) a Document with a valid InstructionSet
func (inst *InstructionSet) CreateDocument() (Document, error) {
	return Document{}, nil
}
