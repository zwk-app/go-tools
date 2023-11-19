package tools

// BoolToString convert a boolean to a string
//
//goland:noinspection GoUnusedExportedFunction
func BoolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
