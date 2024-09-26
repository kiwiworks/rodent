package object

import (
	"net/url"
	"strings"

	"github.com/kiwiworks/rodent/errors"
)

// pathAndKey extracts the first and second elements of the URL path.
// Returns the first element, the second element, and an error if applicable.
func pathAndKey(uri url.URL) (string, string, error) {
	if uri.Scheme != "s3" {
		return "", "", errors.Newf("unsupported scheme %s", uri.Scheme)
	}
	// Remove the leading "/" from the URL path
	path := strings.TrimPrefix(uri.Path, "/")
	if path == "" {
		return "", "", errors.Newf("path is empty")
	}

	// Split the path into a maximum of three segments
	segments := strings.SplitN(path, "/", 3)

	switch len(segments) {
	case 1:
		return segments[0], "", nil
	case 2:
		if segments[1] == "" {
			return "", "", errors.Newf("path is not valid")
		}
		return segments[0], segments[1], nil
	case 3:
		return segments[0], strings.Join(segments[1:], "/"), nil
	default:
		return "", "", errors.Newf("unexpected path format")
	}
}
