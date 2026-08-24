package command

import "testing"

func TestCommand_FullName(t *testing.T) {
	cases := []struct {
		name     string
		command  *Command
		expected string
	}{
		{"top-level", New("alpha", "", ""), "alpha"},
		{"one level deep", New("alpha.run", "", ""), "alpha.run"},
		{"nested", New("a.b.c", "", ""), "a.b.c"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.command.FullName(); got != tc.expected {
				t.Errorf("FullName() = %q, want %q", got, tc.expected)
			}
		})
	}
}
