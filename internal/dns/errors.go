package dns

import "fmt"

type MissingFieldError struct {
	Field string
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("%s is set to the empty value but it is required", e.Field)
}

type UnsupportedOpCodeError struct {
	Opcode int
}

func (e *UnsupportedOpCodeError) Error() string {
	return fmt.Sprintf("opcode %v is not supported at this time", e.Opcode)
}

type UnsupportedRecordTypeError struct {
	Type uint16
}

func (e *UnsupportedRecordTypeError) Error() string {
	return fmt.Sprintf("record type %v is not supported at this time", e.Type)
}
