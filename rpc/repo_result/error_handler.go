package reporesult

import "fmt"

type StoreError struct {
	Code    int
	Message string
}

func (s StoreError) Error() string {
	return fmt.Sprintf("This is MongoError code: %d and message: %s", s.Code, s.Message)
}
