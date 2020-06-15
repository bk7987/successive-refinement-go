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
	booleanArgs     map[rune]bool
	stringArgs      map[rune]string
	intArgs         map[rune]int
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
	a.valid = a.parse()
}

func (a *Args) parse() bool {
	if len(a.schema) == 0 && len(a.args) == 0 {
		return true
	}

	a.parseSchema()
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
	//elementTail := string(element[1:])
	a.validateSchemaElementID(elementID)
}

func (a *Args) validateSchemaElementID(elementID string) {
	if !unicode.IsLetter([]rune(elementID)[0]) {
		message := fmt.Sprintf("Bad character %v in Args format: %v", elementID, a.schema)
		log.Fatal(message)
	}
}

func parseBooleanSchemaElement(elementID rune) {

}

func (a *Args) isBooleanSchemaElement(elementTail string) bool {
	return len(elementTail) == 0
}
