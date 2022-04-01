package models

import (
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
	Evaluate(string) (string, error)
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

func (v stringVariable) Evaluate(value string) (string, error) {
	if value == "" {
		return "___", nil
	}

	return template.HTMLEscapeString(value), nil
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

func (v numberVariable) Evaluate(value string) (string, error) {
	if value == "" {
		return "___", nil
	}

	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}

	if math.Mod(v.Step, 1) == 0 {
		return fmt.Sprintf("%.f", num), nil
	}

	return fmt.Sprintf("%f", num), nil
}

func replaceVariable(verbatimTemplate string, variable AnyVariable, value string) (string, error) {
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
