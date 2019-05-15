package helpers

// DerefString returns a pointer to a variable with the input string. This is useful in models since they use *string.
func DerefString(s string) *string {
	return &s
}
