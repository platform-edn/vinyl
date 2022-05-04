package store

import "fmt"

type MissingRecordError struct {
	Domain string
}

func (e *MissingRecordError) Error() string {
	return fmt.Sprintf("domain %s is not implemented", e.Domain)
}

type ExistingRecordError struct {
	Domain string
}

func (e *ExistingRecordError) Error() string {
	return fmt.Sprintf("a record with the domain %s already exists", e.Domain)
}
