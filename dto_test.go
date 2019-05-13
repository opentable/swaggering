package swaggering

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadPopulate(t *testing.T) {
	assert := assert.New(t)
	target := struct {
		Field string
	}{}
	buf := ioutil.NopCloser(bytes.NewBufferString(`{"Field": "yup"}`))

	err := ReadPopulate(buf, &target)
	assert.NoError(err)

	assert.Equal("yup", target.Field)
}

func TestLoadMap(t *testing.T) {
	assert := assert.New(t)

	m := &TestModel{}
	LoadMap(m, map[string]interface{}{})
	assert.Equal("", m.TestGo)
	LoadMap(m, map[string]interface{}{"TestGo": "yes"})
	assert.Equal("yes", m.TestGo)
}

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	m := &TestModel{}
	j, err := json.Marshal(m)
	assert.NoError(err)
	assert.Equal("{}", string(j))
	m.LoadMap(map[string]interface{}{"TestGo": "test"})
	j, err = json.Marshal(m)
	assert.NoError(err)
	assert.Equal(`{"testSw":"test"}`, string(j))
}

// BEGIN test model object
type TestModel struct {
	present map[string]bool
	TestGo  string `json:"testSw,omitempty"`
}

func (self *TestModel) Populate(jsonReader io.ReadCloser) (err error) {
	return ReadPopulate(jsonReader, self)
}

func (self *TestModel) Absorb(other DTO) error {
	if like, ok := other.(*TestModel); ok {
		*self = *like
		return nil
	}
	return fmt.Errorf("A TestModel cannot absorb the values from %v", other)
}

func (self *TestModel) FormatText() string {
	return FormatText(self)
}

func (self *TestModel) FormatJSON() string {
	return FormatJSON(self)
}

func (self *TestModel) SetField(name string, value interface{}) error {
	if self.present == nil {
		self.present = make(map[string]bool)
	}
	switch name {
	default:
		return fmt.Errorf("No such field %s on TestModel", name)

	case "testSw", "TestGo":
		v, ok := value.(string)
		if ok {
			self.TestGo = v
			self.present["testSw"] = true
			return nil
		} else {
			return fmt.Errorf("Field testSw/TestGo: value %v(%T) couldn't be cast to type string", value, value)
		}

	}
}

func (self *TestModel) GetField(name string) (interface{}, error) {
	switch name {
	default:
		return nil, fmt.Errorf("No such field %s on TestModel", name)

	case "testSw", "TestGo":
		if self.present != nil {
			if _, ok := self.present["testSw"]; ok {
				return self.TestGo, nil
			}
		}
		return nil, fmt.Errorf("Field TestGo no set on TestGo %+v", self)

	}
}

func (self *TestModel) ClearField(name string) error {
	if self.present == nil {
		self.present = make(map[string]bool)
	}
	switch name {
	default:
		return fmt.Errorf("No such field %s on TestModel", name)

	case "testSw", "TestGo":
		self.present["testSw"] = false

	}

	return nil
}

func (self *TestModel) LoadMap(from map[string]interface{}) error {
	return LoadMapIntoDTO(from, self)
}

type TestModelList []*TestModel

func (list *TestModelList) Populate(jsonReader io.ReadCloser) (err error) {
	return ReadPopulate(jsonReader, list)
}

func (list *TestModelList) FormatText() string {
	text := []byte{}
	for _, dto := range *list {
		text = append(text, (*dto).FormatText()...)
		text = append(text, "\n"...)
	}
	return string(text)
}

func (list *TestModelList) FormatJSON() string {
	return FormatJSON(list)
}

// END test object
