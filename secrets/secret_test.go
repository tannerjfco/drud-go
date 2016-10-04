package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretWrite(t *testing.T) {

	args := []string{"secret/test/themagictest1", "name=bob"}

	sobj := Secret{
		Path: NormalizePath(args[0]),
	}

	err := sobj.Init(args)
	assert.NoError(t, err)

	err = sobj.MustValidate()
	assert.NoError(t, err)

	err = sobj.Write()
	assert.NoError(t, err)
}

// TestNormalizePath ensures input paths are translated correctly to what vault expects
func TestIntegrationNormalizePath(t *testing.T) {
	tests := [][]string{
		[]string{"secret", "secret/secret"},
		[]string{"secret/", "secret/secret"},
		[]string{"/", "secret"},
		[]string{"test", "secret/test"},
		[]string{"test/", "secret/test"},
		[]string{"test/tacos", "secret/test/tacos"},
	}
	for _, v := range tests {
		assert.Equal(t, v[1], NormalizePath(v[0]))
	}
}
