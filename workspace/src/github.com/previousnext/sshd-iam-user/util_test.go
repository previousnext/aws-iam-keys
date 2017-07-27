package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReg(t *testing.T) {
	assert.True(t, contains("foo", []string{"foo", "bar"}))
	assert.False(t, contains("baz", []string{"foo", "bar"}))
}
