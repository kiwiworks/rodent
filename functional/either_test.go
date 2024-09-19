package functional

import (
	"testing"
)

type TestStruct struct {
	Field string
}

func TestEither_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name            string
		data            []byte
		expectedErr     bool
		expectedIsLeft  bool
		expectedIsRight bool
	}{
		{
			name:            "Correct JSON left",
			data:            []byte(`{"Field":"Left structured data"}`),
			expectedErr:     false,
			expectedIsLeft:  true,
			expectedIsRight: false,
		},
		{
			name:            "Correct JSON right",
			data:            []byte(`"Right unstructured data"`),
			expectedErr:     false,
			expectedIsLeft:  false,
			expectedIsRight: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := Either[TestStruct, string]{}
			err := e.UnmarshalJSON(tc.data)
			if tc.expectedErr && err == nil {
				t.Errorf("expected an error, but got none")
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("expected no error, but got '%v'", err)
			}
			if e.IsLeft() != tc.expectedIsLeft {
				t.Errorf("expected IsLeft to be '%v', but got '%v'", tc.expectedIsLeft, e.IsLeft())
			}
			if e.IsRight() != tc.expectedIsRight {
				t.Errorf("expected IsRight to be '%v', but got '%v'", tc.expectedIsRight, e.IsRight())
			}
		})
	}
}
