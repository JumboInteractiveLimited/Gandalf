package pathing

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	cases := []struct {
		blob     string
		path     string
		err      bool
		expected []string
	}{
		{``, ``, true, []string{}},
		{`{}`, ``, true, []string{}},
		{
			blob:     `{"name":"test"}`,
			path:     `$.name+`,
			err:      false,
			expected: []string{`"test"`},
		},
		{
			blob:     `{"things":[{"name":"one"},{"name":"two"}]}`,
			path:     `$.things[*]+`,
			err:      false,
			expected: []string{`{"name":"one"}`, `{"name":"two"}`},
		},
		{
			blob:     `{"things":[{"name":"one"},{"name":"two"}]}`,
			path:     `$.things[*].name+`,
			err:      false,
			expected: []string{`"one"`, `"two"`},
		},
		{
			blob:     `{"things":[{"type":"one"},{"type":"two", "unique": 1}]}`,
			path:     `$.things[*]?(@.type == "two").unique+`,
			err:      false,
			expected: []string{`1`},
		},
		{
			blob:     `{}`,
			path:     `$+`,
			err:      false,
			expected: []string{`{}`},
		},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(st *testing.T) {
			result, err := JSON(tt.blob, tt.path)
			testError(st, tt.err, err)
			testStringSlices(st, tt.expected, result)
		})
	}
}

func BenchmarkJSONMultiple(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := JSON(`{"things":[{"name":"one"},{"name":"two"}]}`, `$.things[*].name+`)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONConditional(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := JSON(`{"things":[{"type":"one"},{"type":"two", "unique": 1}]}`,
			`$.things[*]?(@.type == "two").unique+`)
		if err != nil {
			b.Fatal(err)
		}
	}
}
