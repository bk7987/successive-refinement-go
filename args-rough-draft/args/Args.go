package args

import (
	"fmt"
	"log"
	"strings"
	"unicode"
)

// Args defines the structure of the Args class
type Args struct {
	schema          string
	args            []string
	valid           bool
	booleanArgs     map[string]bool
	stringArgs      map[string]string
	intArgs         map[string]int
	currentArgument int
	errorArgumentID rune
	errorParameter  string
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

	}
}

func (a *Args) setArgument(argChar string) bool {
	if a.isBooleanArg(argChar) {
		return true
		//a.setBooleanArg(argChar, true)
	}
	return false
}

func (a *Args) isBooleanArg(argChar string) bool {
	_, ok := a.booleanArgs[argChar]
	fmt.Println(ok)
	return ok
}
