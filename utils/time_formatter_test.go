package utils

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTimeFormatter(t *testing.T) {
	cases := []struct {
		name, in, exp string
	}{
		{
			name: "must cut string",
			in:   "76h29m3123s",
			exp:  "76 ч. 29 м.",
		},
		{
			name: "must cut string",
			in:   "72h29m3123s",
			exp:  "72 ч. 29 м.",
		},
	}

	for i, c := range cases {
		act := CutTimeString(c.in)
		require.Equal(t, c.exp, act, fmt.Sprintf("%d Case: %s", i+1, c.name))
	}

}
