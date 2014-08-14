package lists

type ListErrorType int

const (
	ListErrorOutOfBound ListErrorType = iota
)

type ListError struct {
	Type ListErrorType
}

func (err *ListError) Error() string {
	switch err.Type {
	default:
		return "Out of bound"
	}
}
