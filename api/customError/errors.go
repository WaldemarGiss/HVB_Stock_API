package customError

import "fmt"

type ErrorStock struct {
	Code int
	Text string
}

func (err *ErrorStock) Error() string {
	return fmt.Sprintf("stock-api returned status : %v \n and text : %v", err.Code, err.Text)
}
