package dto

import "fmt"

func ErrMissingField(field string) error {
	return fmt.Errorf("%s is required", field)
}
