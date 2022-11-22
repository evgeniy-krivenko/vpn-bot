package handler

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEscape(t *testing.T) {
	cases := []struct {
		name string
		in   string
		exp  string
	}{
		{
			name: "one_dot",
			in:   "with dot.",
			exp:  "with dot\\.",
		},
		{
			name: "many_dots",
			in:   "with.. dot.",
			exp:  "with\\.\\. dot\\.",
		},
		{
			name: "other_symbols",
			in:   "()[]",
			exp:  "\\(\\)\\[\\]",
		},
		{
			name: "back quotes",
			in:   "`some-text`",
			exp:  "\\`some\\-text\\`",
		},
	}

	for _, c := range cases {
		resp := escape(c.in)
		require.Equal(t, c.exp, resp)
	}
}
