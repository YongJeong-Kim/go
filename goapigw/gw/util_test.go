package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMap2Struct(t *testing.T) {
	type m map[string]string
	type param struct {
		pm m
	}
	var s struct {
		S1 string
		S2 string
	}

	testCases := []struct {
		name  string
		p     param
		check func(p param)
	}{
		{
			name: "OK",
			p: param{
				pm: m{
					"s1": "asdf",
					"s2": "fdas",
				},
			},
			check: func(p param) {
				result, err := Map2Struct(&s, p.pm)
				require.NoError(t, err)

				require.Equal(t, p.pm["s1"], result.S1)
				require.Equal(t, p.pm["s2"], result.S2)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.check(tc.p)
		})
	}
}
