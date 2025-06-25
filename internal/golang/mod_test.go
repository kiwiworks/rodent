package golang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParentModFromPath(t *testing.T) {
	r := require.New(t)

	mod, err := FindModulePath("./internal/golang/mod_test.go")
	r.NoError(err)
	r.Equal("github.com/kiwiworks/rodent", mod)
}

func TestFindModulePath_Errors(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		expectedError bool
	}{
		{
			name:          "Empty path",
			path:          "",
			expectedError: true,
		},
		{
			name:          "Non-existent path",
			path:          "/completely/nonexistent/path",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			mod, err := FindModulePath(tc.path)
			if tc.expectedError {
				r.Error(err)
				r.Empty(mod)
			} else {
				r.NoError(err)
				r.NotEmpty(mod)
			}
		})
	}
}
