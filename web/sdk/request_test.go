package sdk

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithQueryParam(t *testing.T) {
	testCases := []struct {
		name  string
		param string
		value string
	}{
		{
			name:  "single param",
			param: "param1",
			value: "value1",
		},
		{
			name:  "empty param",
			param: "",
			value: "value1",
		},
		{
			name:  "empty value",
			param: "param1",
			value: "",
		},
		{
			name:  "empty param & value",
			param: "",
			value: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init request with empty values
			request := &Request{
				Values: make(url.Values),
			}

			// apply test param and value
			WithQueryParam(tc.param, tc.value)(request)

			require.Equal(t, 1, len(request.Values[tc.param]), "unexpected number of values for param")
			require.Equal(t, tc.value, request.Values[tc.param][0], "unexpected param value")
		})
	}
}

func TestWithMultipleQueryParams(t *testing.T) {
	request := &Request{
		Values: make(url.Values),
	}

	// apply multiple params and values
	WithQueryParam("param1", "value1")(request)
	WithQueryParam("param1", "value2")(request)
	WithQueryParam("param2", "value1")(request)

	require.Equal(t, 2, len(request.Values["param1"]), "unexpected number of values for param1")
	require.Equal(t, "value1", request.Values["param1"][0], "unexpected 1st value for param1")
	require.Equal(t, "value2", request.Values["param1"][1], "unexpected 2nd value for param1")

	require.Equal(t, 1, len(request.Values["param2"]), "unexpected number of values for param2")
	require.Equal(t, "value1", request.Values["param2"][0], "unexpected 1st value for param2")
}
