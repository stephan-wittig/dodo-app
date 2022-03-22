package models

// Instructive is the interface that should be implemented by all variables
// it makes sure that they can be turned into serialisable intsructions
type Instructive interface {
	// ToInstructions should return key and instruction
	ToInstructions() variableInstruction
	GetKey() string
}

// variable composes the common attributes for all variables
type variable struct {
	Key         string `xml:"key,attr"`
	Label       string `xml:"label,attr"`
	Description string ` xml:"description,attr"` // Optional
}

type stringVariable struct {
	MaxLength float64 `xml:"maxLength,attr"` // Default so 280 (a Tweet)
	MinLength float64 `xml:"minLength,attr"` // Defaults to 0

	variable
}

type numberVariable struct {
	Max  float64 `xml:"max,attr"`  // Optional
	Min  float64 `xml:"min,attr"`  // Defaults to 0
	Step float64 `xml:"step,attr"` // Defaults to 1. Step size, eg. 0.01 for currencies

	variable
}

// choiceVariable represents the choice between multiple options (default) or a list of options
// that can be toggled, depending on MaxChoices
type choiceVariable struct {
	MaxChoices float64        `xml:"maxChoices,attr"` // Optional
	MinChoices float64        `xml:"minChoices,attr"` // Defaults to 1
	Options    []choiceOption `xml:"Option"`

	variable
}

type choiceOption struct {
	Label string `json:"label,omitempty" xml:"Label,attr"`
	Value string `json:"value" xml:",chardata"`
}

// Methods

func (v stringVariable) ToInstructions() variableInstruction {
	return variableInstruction{
		DataType:    "STRING",
		Description: v.Description,
		Min:         v.MinLength,
		Max:         v.MaxLength,
	}
}

func (v stringVariable) GetKey() string {
	return v.Key
}

func (v numberVariable) ToInstructions() variableInstruction {
	return variableInstruction{
		DataType:    "NUMBER",
		Description: v.Description,
		Min:         v.Min,
		Max:         v.Max,
		Step:        v.Step,
	}
}

func (v numberVariable) GetKey() string {
	return v.Key
}

func (v choiceVariable) ToInstructions() variableInstruction {
	return variableInstruction{
		DataType:    "CHOICE",
		Description: v.Description,
		Min:         v.MinChoices,
		Max:         v.MaxChoices,
		Options:     v.Options,
	}
}

func (v choiceVariable) GetKey() string {
	return v.Key
}
