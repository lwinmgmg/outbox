package helper

import "fmt"

func PanicEmptyString(input any, message ...string) {
	var mesg string
	if len(message) == 0 {
		mesg = "Missing one required field"
	} else {
		mesg = fmt.Sprintf("%v is required field", message[0])
	}
	switch input.(type) {
	case int:
		if input == 0 {
			panic(mesg)
		}
	case string:
		if input == "" {
			panic(mesg)
		}
	}
}

func DefaultString(input string, value string) string {
	if input == "" {
		return value
	}
	return input
}

func DefaultInt(input int, value int) int {
	if input == 0 {
		return value
	}
	return input
}

func DefaultBoolean(input bool, value bool) bool {
	if !input {
		return value
	}
	return input
}
