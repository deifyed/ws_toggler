package filesystem

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeal(t *testing.T) {
	testCases := []struct {
		name string

		withContent  []byte
		expectResult []byte
	}{
		{
			name: "Should be safe to use on correct result",

			withContent:  []byte("{\"previous_workspace\":\"2\"}"),
			expectResult: []byte("{\"previous_workspace\":\"2\"}"),
		},
		{
			name: "Should remove duplicate curly bracket",

			withContent:  []byte("{\"previous_workspace\":\"2\"}}"),
			expectResult: []byte("{\"previous_workspace\":\"2\"}"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result := heal(tc.withContent)

			assert.True(t, bytes.Equal(tc.expectResult, result))
		})
	}
}
