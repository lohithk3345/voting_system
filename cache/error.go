package cache

import "fmt"

const (
	SET_ERROR = iota
	GET_ERROR
	SET_MARSHAL_ERROR
)

type CacheError struct {
	Code    int
	Message string
}

func (c CacheError) Error() string {
	return fmt.Sprintf("This is CacheService Error Code: %d and Message: %s", c.Code, c.Message)
}
