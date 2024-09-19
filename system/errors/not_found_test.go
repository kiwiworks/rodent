package errors

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "NotFound error",
			err:  Newf("another error"),
			want: false,
		},
		{
			name: "Wrapped NotFoundError",
			err:  errors.Wrapf(NotFound("test entity", uuid.New()), "additional context"),
			want: true,
		},
		{
			name: "Not a NotFoundError",
			err:  fmt.Errorf("another error"),
			want: false,
		},
		{
			name: "No error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
