package unpack

// HasSupportedTool returns true if at least one tool exists in PATH.
func HasSupportedTool() bool {
	_, err := DetectTool(nil)
	return err == nil
}
