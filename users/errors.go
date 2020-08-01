package users

// Error represents errors related to users
// They should always be checked using errors.Is since the
// error may be wrapped with more context.
type Error string

const (
	InvalidUserParamErr Error = "user has invalid param"
)

// Error returns the string representation of the error
func (e Error) Error() string {
	return string(e)
}
