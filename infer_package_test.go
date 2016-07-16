package swaggering

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInferPackage(t *testing.T) {
	assert := assert.New(t)

	p, err := packageUnderGopath("/gopath", "/gopath/src/github.com/opentable/otthing")
	assert.NoError(err)
	assert.Equal("github.com/opentable/otthing", p)

	p, err = packageUnderGopath("/gopath", "/somwhere/completly/other")
	assert.Error(err)
}
