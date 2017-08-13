package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunExec(t *testing.T) {

	assert := assert.New(t)

	ls, err := Run("./", "ls -la")
	assert.Nil(err, "Run returned an unexpected error")
	if assert.NotNil(ls, "Run ls -la failed to return a result") {
		assert.Contains(string(ls), "osutil.go", "Run 'ls -la' failed to contain file in chosen folder")
		assert.Contains(string(ls), "osutil_test.go", "Run 'ls -la' failed to contain second file in chosen folder")
	}

	// Test with many paramaters
	ls, err = Run("../../../", `stat -f %N-%A main.go`)
	assert.Nil(err, "Run returned an unexpected error")
	if assert.NotNil(ls, "Run ls -la failed to return a result") {
		assert.Contains(string(ls), "main.go-644", "Run 'stat -f' failed to contain file in chosen folder")
	}
}
