package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIP(t *testing.T) {
	testCases := []struct {
		raw, parsed string
	}{
		{"127.0.0.1", "127.0.0.1"},
		{"10.0.0.1", "10.0.0.1"},
		{"10.0.0.1:1234", "10.0.0.1"},
		{"::1", "::1"},
		{"[::1]:1234", "::1"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.raw, func(t *testing.T) {
			assert.Equal(t, testCase.parsed, parseIP(testCase.raw))
		})
	}
}
