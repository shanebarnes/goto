package aggregator

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregator_Insert(t *testing.T) {
	a := NewAggregator()

	assert.Equal(t, syscall.ENOTSUP, a.Insert("int", "0"))
	assert.Nil(t, a.Insert("int", int(1)))
	assert.Nil(t, a.Insert("int", int(2)))
	assert.Equal(t, syscall.EINVAL, a.Insert("int", int8(3)))
	assert.Equal(t, syscall.EINVAL, a.Insert("int", "4"))
}

func TestAggregator_NewAggregator(t *testing.T) {
	a := NewAggregator()

	assert.NotNil(t, a)
	assert.NotNil(t, a.db)
}
