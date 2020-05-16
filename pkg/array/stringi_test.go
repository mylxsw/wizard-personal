package array_test

import (
	"github.com/mylxsw/wizard-personal/pkg/array"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringUnique(t *testing.T) {
	arr := []string{
		"aaa",
		"bbb",
		"ccc",
		"aaa",
		"ddd",
		"ccc",
	}

	assert.EqualValues(t, 4, len(array.StringUnique(arr)))
}

