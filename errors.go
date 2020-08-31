package configo

type ErrInvalidSpec struct{}

func (e *ErrInvalidSpec) Error() string {
	return "spec isn't valid"
}
