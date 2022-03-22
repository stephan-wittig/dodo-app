package models

// Document is the result of instanciating a template with a set of variables
//
// The actual document ("verbatim") as well as the instructions are included.
type Document struct {
	Id     string           `json:"id,omitempty"`     // For stored documents only; UUID
	Digest []byte           `json:"digest"`           // MD5 of InstructionSet
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
