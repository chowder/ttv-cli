package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemove(t *testing.T) {
	lst := []string{"a", "b", "c"}
	lst = Remove(lst, "b")

	assert.Equal(t, lst, []string{"a", "c"})
}
