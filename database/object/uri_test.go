package object

import (
	"net/url"
	"testing"
)

func TestPathAndKey(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantPath string
		wantKey  string
		wantErr  bool
	}{
		{
			name:     "valid path with key",
			url:      "s3://example.com/somepath/somekey",
			wantPath: "somepath",
			wantKey:  "somekey",
			wantErr:  false,
		},
		{
			name:     "valid path without key",
			url:      "s3://example.com/somepath",
			wantPath: "somepath",
			wantKey:  "",
			wantErr:  false,
		},
		{
			name:     "invalid path without key and /",
			url:      "s3://example.com",
			wantPath: "",
			wantKey:  "",
			wantErr:  true,
		},
		{
			name:     "invalid path with / at the end",
			url:      "s3://example.com/somepath/",
			wantPath: "",
			wantKey:  "",
			wantErr:  true,
		},
		{
			name:     "valid path with multiple / and key",
			url:      "s3://example.com/somepath/anotherpath/key/foo",
			wantPath: "somepath",
			wantKey:  "anotherpath/key/foo",
			wantErr:  false,
		},
		{
			name:     "invalid scheme",
			url:      "https://example.com/somepath/valid",
			wantPath: "",
			wantKey:  "",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := url.Parse(test.url)
			if err != nil {
				t.Fatal(err)
			}

			gotPath, gotKey, err := pathAndKey(*uri)
			if (err != nil) != test.wantErr {
				t.Errorf("pathAndKey() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotPath != test.wantPath {
				t.Errorf("pathAndKey() gotPath = %v, want %v", gotPath, test.wantPath)
			}
			if gotKey != test.wantKey {
				t.Errorf("pathAndKey() gotKey = %v, want %v", gotKey, test.wantKey)
			}
		})
	}
}
