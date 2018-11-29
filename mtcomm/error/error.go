package error

type BusinessError struct {
	error
	Code string
	Msg  string
}

func NewBussinessError(code string, msg string) error {
	return BusinessError{
		Code: code,
		Msg:  msg,
	}
}

func (e BusinessError) Error() string {
	return e.Code + ": " + e.Msg
}
