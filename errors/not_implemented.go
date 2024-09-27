package errors

func NotImplemented(message ...string) error {
	reason := "this feature is not implemented yet"
	if len(message) > 0 {
		reason = message[0]
	}
	return Newf("not implemented: %s", reason)
}
