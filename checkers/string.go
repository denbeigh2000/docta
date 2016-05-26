package checkers

import (
	"github.com/denbeigh2000/docta"

	"strings"
)

// CheckString is a utility function for checking if a given string contains
// a substring that triggers an error state.
func CheckString(given string, yellow string, red string) docta.State {
	if strings.Contains(given, red) {
		return docta.Red
	}

	if strings.Contains(given, yellow) {
		return docta.Yellow
	}

	return docta.Green
}
