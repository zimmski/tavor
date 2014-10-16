package lists

// ListErrorType the list error type
type ListErrorType int

const (
	// ListErrorOutOfBound an index not in the bound of available list items was used.
	ListErrorOutOfBound ListErrorType = iota
)

// ListError holds a list error
type ListError struct {
	Type ListErrorType
}

func (err *ListError) Error() string {
	switch err.Type {
	default:
		return "Out of bound"
	}
}
