package swaggering

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeedDTO(t *testing.T) {
	assert := assert.New(t)

	clnt, ctrl := NewChannelDummy()

	m := &TestModel{}
	m.SetField("TestGo", "henry")
	ctrl.FeedDTO(m, nil)

	other := &TestModel{}
	prms := UrlParams{}
	err := clnt.DTORequest("test-somewhere", other, "GET", "/somewhere", prms, prms)
	assert.Equal(m.TestGo, other.TestGo)
	ov, err := other.GetField("TestGo")
	assert.NoError(err)
	assert.Equal("henry", ov)
}

func TestFeedSimple(t *testing.T) {
	assert := assert.New(t)

	clnt, ctrl := NewChannelDummy()

	ctrl.FeedSimple("hi", nil)

	prms := UrlParams{}
	body, err := clnt.Request("test-whatever", "GET", "/whatever", prms, prms)
	assert.NoError(err)
	got, err := ioutil.ReadAll(body)
	assert.NoError(err)
	assert.Equal("hi", string(got))
}
