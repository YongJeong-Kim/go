package main

import (
	"github.com/stretchr/testify/require"
	"goregex/example"
	"testing"
)

func TestNumberOnly(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		check func(input string)
	}{
		{
			name:  "num",
			input: "235235",
			check: func(input string) {
				ok, err := NumberOnly(input)
				require.NoError(t, err)
				require.True(t, ok)
			},
		},
		{
			name:  "string",
			input: "dfse",
			check: func(input string) {
				ok, err := NumberOnly(input)
				require.NoError(t, err)
				require.False(t, ok)
			},
		},
		{
			name:  "string and num",
			input: "qwe22rqwe4",
			check: func(input string) {
				ok, err := NumberOnly(input)
				require.NoError(t, err)
				require.False(t, ok)
			},
		},
		{
			name:  "string, num, spec",
			input: "@eq21$fe1Q3",
			check: func(input string) {
				ok, err := NumberOnly(input)
				require.NoError(t, err)
				require.False(t, ok)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(tc.input)
		})
	}
}

func TestSituation1(t *testing.T) {
	example.Password10()
}
