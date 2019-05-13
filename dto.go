package swaggering

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type (
	// DTO is an interface for a generic data transfer object.
	DTO interface {
		Populate(io.ReadCloser) error
		Absorb(DTO) error
		FormatText() string
		FormatJSON() string
	}

	// Fielder is an interface for an object with optional fields
	// This is the most surprising aspect of swaggering, but obvious on reflection.
	// JSON interfaces often treat the _absence_ of a field as very different from
	// its presence, regardless of the value of the field. { name: "Judson" } is
	// semantically different from { name: "Judson", job: undefined }.
	// It's important to distinguish absence from zero, therefore.
	Fielder interface {
		GetField(string) (interface{}, error)
		SetField(string, interface{}) error
		ClearField(string) error
		LoadMap(map[string]interface{}) error
	}
)

// ReadPopulate reads from jsonReader in order to fill in target
func ReadPopulate(jsonReader io.ReadCloser, target interface{}) error {
	dec := json.NewDecoder(jsonReader)
	return dec.Decode(target)
}

// LoadMap loads a map of values into a Fielder
func LoadMap(dto Fielder, from map[string]interface{}) (Fielder, error) {
  return dto, LoadMapIntoDTO(from, dto)
}

// FormatText formats a DTO
func FormatText(dto interface{}) string {
	return fmt.Sprintf("%+v", dto)
}

// FormatJSON formats a dto as JSON
func FormatJSON(dto interface{}) string {
	str, err := json.Marshal(dto)
	if err != nil {
		return "&lt;<XXXX>>"
	}
	buf := bytes.Buffer{}
	json.Indent(&buf, str, "", "  ")
	return buf.String()
}

// LoadMapIntoDTO loads a map of key/values into a DTO, setting their presence
// as they're loaded
func LoadMapIntoDTO(from map[string]interface{}, dto Fielder) error {
	errs := make([]string, 0)
	for name, value := range from {
		if err := dto.SetField(name, value); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// vim: set ft=go:
