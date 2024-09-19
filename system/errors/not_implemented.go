package errors

func NotImplemented() error {
	return Newf(
		"not implemented: %s",
		"this feature is not implemented yet")
}
