package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"text/template"
)

// globalVariables holds the global variables. This should be retrieved from somewhere.
// HTML should already be escaped
var globalVariables = map[string]string{
	"company-full-name":  "Shiny App LLC",
	"company-short-name": "Shiny App",
	"active-from":        "April 1st 2022",
}

// olThreshhold is a list with choices that are longer than this will be rendered as ol
var olThreshhold = 30

// AnyVariable is the interface that should be implemented by all variables
// it makes sure that they can be turned into serialisable intsructions
type AnyVariable interface {
	ToInstructions() variableInstruction
	GetKey() string
	// Evaluates turns the input data into a string fitting the format. HTML is escaped here
	Evaluate([]byte) (string, error)
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
		Label:       v.Label,
		Description: v.Description,
		Min:         v.MinLength,
		Max:         v.MaxLength,
	}
}

func (v stringVariable) GetKey() string {
	return v.Key
}

func (v stringVariable) Evaluate(value []byte) (string, error) {
	if value == nil {
		return "___", nil
	}

	return template.HTMLEscapeString(string(value)), nil
}

func (v numberVariable) ToInstructions() variableInstruction {
	return variableInstruction{
		DataType:    "NUMBER",
		Label:       v.Label,
		Description: v.Description,
		Min:         v.Min,
		Max:         v.Max,
		Step:        v.Step,
	}
}

func (v numberVariable) GetKey() string {
	return v.Key
}

func (v numberVariable) Evaluate(value []byte) (string, error) {
	if value == nil {
		return "___", nil
	}

	num, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return "", err
	}

	if math.Mod(v.Step, 1) == 0 {
		return fmt.Sprintf("%.f", num), nil
	}

	return fmt.Sprintf("%f", num), nil
}

func (v choiceVariable) ToInstructions() variableInstruction {
	return variableInstruction{
		DataType:    "CHOICE",
		Label:       v.Label,
		Description: v.Description,
		Min:         v.MinChoices,
		Max:         v.MaxChoices,
		Options:     v.Options,
	}
}

func (v choiceVariable) GetKey() string {
	return v.Key
}

func (v choiceVariable) Evaluate(value []byte) (string, error) {
	if value == nil {
		return "___", nil
	}

	var choices []int64
	if err := json.Unmarshal(value, &choices); err != nil {
		return "", err
	}

	choicesVerbatim := []string{}
	for _, choice := range choices {
		choicesVerbatim = append(choicesVerbatim, v.Options[choice].Value)
	}

	sliceLen := len(choices)

	if sliceLen == 1 {
		return choicesVerbatim[0], nil
	}

	maxChoiceLen := 0

	for _, choice := range choicesVerbatim {
		choiceLen := len(choice)
		if choiceLen > maxChoiceLen {
			maxChoiceLen = choiceLen
		}
	}

	items := ""

	if maxChoiceLen > olThreshhold {
		for _, choice := range choicesVerbatim {
			items = items + fmt.Sprintf(
				"<li>%s</li>\n",
				choice,
			)
		}
		return fmt.Sprintf("\n<ol>\n%s</ol>\n", items), nil
	}

	for i, choice := range choicesVerbatim {
		items = items + fmt.Sprintf(
			"(%d)\u00A0%s ", // 00A0 = non-breaking space
			i+1,
			choice,
		)
	}
	return strings.TrimSpace(items), nil
}

func replaceVariable(verbatimTemplate string, variable AnyVariable, value []byte) (string, error) {
	new, err := variable.Evaluate(value)
	if err != nil {
		return "", err
	}

	return strings.Replace(
		verbatimTemplate,
		fmt.Sprintf("${%s}", variable.GetKey()),
		string(new),
		-1,
	), err
}

func replaceGlobalVariables(verbatimTemplate string) string {
	copy := verbatimTemplate
	for k, v := range globalVariables {
		copy = strings.Replace(
			copy,
			fmt.Sprintf("${{%s}}", k),
			v,
			-1,
		)
	}
	return copy
}
