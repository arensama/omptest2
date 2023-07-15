package errWithStatus

type StatusErr struct {
	Status  int
	Message string
	Err     error
}

func (se StatusErr) Error() string {
	return se.Message
}
func (se StatusErr) Unwrap() error {
	return se.Err
}
