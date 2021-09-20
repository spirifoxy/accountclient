package accounts

import "fmt"

// RequestStatusError is used when response status code is
// not what we expect to receive
type RequestStatusError struct {
	Err  error
	Code int
}

func (e *RequestStatusError) Error() string {
	return fmt.Sprintf("request error, status code %d. %v", e.Code, e.Err)
}

// ClientInitError is used when something unexpected happens while
// creating the api client
type ClientInitError struct {
	Err error
}

func (e *ClientInitError) Error() string {
	return fmt.Sprintf("client initialization error. %v", e.Err)
}
