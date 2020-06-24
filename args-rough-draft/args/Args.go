package args

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// Args defines the structure of the Args class
type Args struct {
	schema              string
	args                []string
	valid               bool
	unexpectedArguments map[string]string
	booleanArgs         map[string]bool
	stringArgs          map[string]string
	intArgs             map[string]int
	argsFound           map[string]string
	currentArgument     int
	errorArgumentID     string
	errorParameter      string
	errorCode           string
}

// Defines all possible error codes
const (
	ErrorCodeOk                 = "OK"
	ErrorCodeMissingString      = "MISSING_STRING"
	ErrorCodeMissingInteger     = "MISSING_INTEGER"
	ErrorCodeInvalidInteger     = "INVALID_INTEGER"
	ErrorCodeUnexpectedArgument = "UNEXPECTED_ARGUMENT"
)

// Init initializes the Args instance
func (a *Args) Init(schema string, args []string) {
	a.schema = schema
	a.args = args

	a.booleanArgs = map[string]bool{}
	a.stringArgs = map[string]string{}
	a.intArgs = map[string]int{}
	a.argsFound = map[string]string{}
	a.unexpectedArguments = map[string]string{}
	a.errorCode = ErrorCodeOk

	a.valid = a.parse()
}

func (a *Args) parse() bool {
	if len(a.schema) == 0 && len(a.args) == 0 {
		return true
	}

	a.parseSchema()
	a.parseArguments()
	return a.valid
}

func (a *Args) parseSchema() bool {
	for _, element := range strings.Split(a.schema, ",") {
		if len(element) > 0 {
			trimmedElement := strings.Trim(element, " ")
			a.parseSchemaElement(trimmedElement)
		}
	}

	return true
}

func (a *Args) parseSchemaElement(element string) {
	elementID := string(element[0])
	elementTail := string(element[1:])
	a.validateSchemaElementID(elementID)

	if isBooleanSchemaElement(elementTail) {
		a.parseBooleanSchemaElement(elementID)
	} else if isStringSchemaElement(elementTail) {
		a.parseStringSchemaElement(elementID)
	} else if isIntegerSchemaElement(elementTail) {
		a.parseIntegerElementSchema(elementID)
	} else {
		message := fmt.Sprintf("Argument %v has invalid format: %v", elementID, elementTail)
		log.Fatal(message)
	}
}

func (a *Args) validateSchemaElementID(elementID string) {
	if !unicode.IsLetter([]rune(elementID)[0]) {
		message := fmt.Sprintf("Bad character %v in Args format: %v", elementID, a.schema)
		log.Fatal(message)
	}
}

func (a *Args) parseBooleanSchemaElement(elementID string) {
	a.booleanArgs[elementID] = false
}

func (a *Args) parseIntegerElementSchema(elementID string) {
	a.intArgs[elementID] = 0
}

func (a *Args) parseStringSchemaElement(elementID string) {
	a.stringArgs[elementID] = ""
}

func isStringSchemaElement(elementTail string) bool {
	return elementTail == "*"
}

func isBooleanSchemaElement(elementTail string) bool {
	return len(elementTail) == 0
}

func isIntegerSchemaElement(elementTail string) bool {
	return elementTail == "#"
}

func (a *Args) parseArguments() bool {
	for a.currentArgument = 0; a.currentArgument < len(a.args); a.currentArgument++ {
		arg := a.args[a.currentArgument]
		a.parseArgument(arg)
	}

	return true
}

func (a *Args) parseArgument(arg string) {
	if string(arg[0]) == "-" {
		a.parseElements(arg)
	}
}

func (a *Args) parseElements(arg string) {
	for _, element := range arg {
		a.parseElement(string(element))
	}
}

func (a *Args) parseElement(argChar string) {
	if a.setArgument(argChar) {
		a.argsFound[argChar] = argChar
	} else {
		a.unexpectedArguments[argChar] = argChar
		a.errorCode = ErrorCodeUnexpectedArgument
		a.valid = false
	}
}

func (a *Args) setArgument(argChar string) bool {
	if a.isBooleanArg(argChar) {
		a.setBooleanArg(argChar, true)
	} else if a.isStringArg(argChar) {
		a.setStringArg(argChar)
	} else if a.isIntArg(argChar) {
		a.setIntArg(argChar)
	} else {
		return false
	}
	return true
}

func (a *Args) isIntArg(argChar string) bool {
	_, ok := a.intArgs[argChar]
	return ok
}

func (a *Args) setIntArg(argChar string) {
	a.currentArgument++
	parameter := a.args[a.currentArgument]
	intParam, err := strconv.Atoi(parameter)

	if err != nil {
		a.errorArgumentID = argChar
		a.errorParameter = parameter
		a.errorCode = ErrorCodeInvalidInteger
	}

	a.intArgs[argChar] = intParam
}

func (a *Args) setStringArg(argChar string) {
	a.currentArgument++
	a.stringArgs[argChar] = a.args[a.currentArgument]
}

func (a *Args) isStringArg(argChar string) bool {
	_, ok := a.stringArgs[argChar]
	return ok
}

func (a *Args) setBooleanArg(argChar string, value bool) {
	a.booleanArgs[argChar] = value
}

func (a *Args) isBooleanArg(argChar string) bool {
	_, ok := a.booleanArgs[argChar]
	return ok
}

// Cardinality returns the number of args found
func (a *Args) Cardinality() int {
	return len(a.argsFound)
}

// Usage returns the current schema for the Args instance
func (a *Args) Usage() string {
	if len(a.schema) > 0 {
		return fmt.Sprintf("-[%v]", a.schema)
	}

	return ""
}

// ErrorMessage returns the current error message
func (a *Args) ErrorMessage() string {
	switch a.errorCode {
	case ErrorCodeOk:
		log.Fatal("TILT: should not get here")
	case ErrorCodeUnexpectedArgument:
		return a.unexpectedArgumentMessage()
	case ErrorCodeMissingString:
		return fmt.Sprintf("Could not find string parameter for -%v", a.errorArgumentID)
	case ErrorCodeInvalidInteger:
		return fmt.Sprintf("Argument -%v expects an integer but was '%v'.", a.errorArgumentID, a.errorParameter)
	case ErrorCodeMissingInteger:
		return fmt.Sprintf("Could not find integer parameter for -%v", a.errorArgumentID)
	}

	return ""
}

func (a *Args) unexpectedArgumentMessage() string {
	message := "Argument(s) -"
	for arg := range a.unexpectedArguments {
		message += arg
	}
	message += " unexpected."
	return message
}

func falseIfNull(b bool) bool {
	return b
}

func zeroIfNull(i int) int {
	return i
}

func blankIfNull(s string) string {
	return s
}

// GetString returns the requested string argument
func (a *Args) GetString(arg string) string {
	return blankIfNull(a.stringArgs[arg])
}

// GetInt returns the requested integer argument
func (a *Args) GetInt(arg string) int {
	return zeroIfNull(a.intArgs[arg])
}

// GetBoolean returns the requested boolean argument
func (a *Args) GetBoolean(arg string) bool {
	return falseIfNull(a.booleanArgs[arg])
}

// Has returns true if the requested argument has been found
func (a *Args) Has(arg string) bool {
	_, has := a.argsFound[arg]
	return has
}

// IsValid returns the valid variable
func (a *Args) IsValid() bool {
	return a.valid
}
