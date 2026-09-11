package httpport

type Validatable interface {
	Validate() error
}
