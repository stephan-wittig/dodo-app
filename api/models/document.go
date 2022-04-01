package models

import (
	"bytes"
	"text/template"

	"github.com/stephan-wittig/dodo/utils"
)

// Document is the result of instanciating a template with a set of variables
//
// The actual document ("verbatim") as well as the instructions are included.
type Document struct {
	Id     string           `json:"id,omitempty"`     // For stored documents only; UUID
	Digest [16]byte         `json:"digest"`           // MD5 of InstructionSet
	Valid  bool             `json:"valid"`            // If false, see errors for details
	Errors map[string]error `json:"errors,omitempty"` // Map of key => error

	// Basic html of document
	//
	// The following html elements are allowed:
	//
	// * section, p, h1, h2
	//
	// * ul, ol, il
	//
	// * em, strong
	//
	// Root should be html, followed by body
	Verbatim []byte `json:"verbatim,omitempty"`

	InstructionSet
}

// intermediateDocument is an intermediate step between template and document
//
// Variables are already evaluated here
type intermediateDocument struct {
	Name     string
	Preambel string
	Sections []intermediateSection
}

// intermediatSection is an intermediate step between sectionTeomplate and document
//
// Variables are already evaluated here
type intermediateSection struct {
	Heading    string
	Preambel   string
	Paragraphs []string
}

func (doc *intermediateDocument) GenerateVerbatim() ([]byte, error) {
	tmplText, err := utils.OpenFile("../demo", "*.template")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("html").Parse(string(tmplText))
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, doc); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
