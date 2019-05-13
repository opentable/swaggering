package swaggering

import (
	"fmt"
	"io"
	"strings"
)

// StringList - it's a list, now with 100% more string
// This type is used in generated DTOs to encode lists of strings.
type StringList []string

// Populate loads a StringList from json
func (list *StringList) Populate(jsonReader io.ReadCloser) (err error) {
	return ReadPopulate(jsonReader, list)
}

// FormatText formats a StringList as text
func (list *StringList) FormatText() string {
	return strings.Join(*list, "\n")
}

// FormatJSON formats a StringList to JSON
func (list *StringList) FormatJSON() string {
	return FormatJSON(list)
}

// Absorb implements DTO for StringList
func (list *StringList) Absorb(other DTO) error {
	if like, ok := other.(*StringList); ok {
		*list = *like
		return nil
	}

	return fmt.Errorf("A StringList cannot absorb from a %T (%v)", other, other)

}
