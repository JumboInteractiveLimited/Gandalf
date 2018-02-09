package pathing

import (
	"fmt"
	"testing"
)

func TestGJSON(t *testing.T) {
	cases := []struct {
		blob     string
		path     string
		err      bool
		expected []string
	}{
		{``, ``, false, []string{}},
		{`{}`, ``, false, []string{"{}"}},
		{
			blob:     `{"name":"test"}`,
			path:     `name`,
			err:      false,
			expected: []string{`"test"`},
		},
		{
			blob:     `{"name":"test"}`,
			path:     `names`,
			err:      true,
			expected: []string{},
		},
		{
			blob:     `{"things":[{"name":"one"},{"name":"two"}]}`,
			path:     `things`,
			err:      false,
			expected: []string{`{"name":"one"}`, `{"name":"two"}`},
		},
		{
			blob:     `{"things":{"ids":[]}`,
			path:     `things`,
			err:      false,
			expected: []string{`{"ids":[]}`},
		},
		{
			blob:     `{"things":[{"name":"one"},{"name":"two"}]}`,
			path:     `things.#.name`,
			err:      false,
			expected: []string{`"one"`, `"two"`},
		},
		{
			blob:     `{"things":[{"name":"one"},{"name":"two"}]}`,
			path:     `things.#.names`,
			err:      true,
			expected: []string{},
		},
		{
			blob:     `{"things":[{"type":"one"},{"type":"two", "unique": 1}]}`,
			path:     `things.#[type=="two"]#.unique`,
			err:      false,
			expected: []string{`1`},
		},
		{
			blob:     `{"things":[{"type":"one"},{"type":"two", "unique": 1}]}`,
			path:     `things.#[type=="two"]#`,
			err:      false,
			expected: []string{`{"type":"two", "unique": 1}`},
		},
		{
			blob:     `{}`,
			path:     ``,
			err:      false,
			expected: []string{`{}`},
		},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(st *testing.T) {
			result, err := GJSON(tt.blob, tt.path)
			testError(st, tt.err, err)
			testStringSlices(st, tt.expected, result)
		})
	}
}

func BenchmarkGJSONMultiple(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := GJSON(`{"things":[{"name":"one"},{"name":"two"}]}`, `things.#.name`)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGJSONConditional(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := GJSON(`{"things":[{"type":"one"},{"type":"two", "unique": 1}]}`,
			`things.#[type=="two"]#.unique`)
		if err != nil {
			b.Fatal(err)
		}
	}
}
