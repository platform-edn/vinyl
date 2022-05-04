package vinyl

import (
	"fmt"
	"net"

	valid "github.com/asaskevich/govalidator"
)

type Record struct {
	Domain  string
	Address string
	TTL     uint32
}

type InvalidRecordAddressError struct {
	Address string
}

func (e *InvalidRecordAddressError) Error() string {
	return fmt.Sprintf("%s is set to the empty value but it is required", e.Address)
}

type InvalidRecordDomainError struct {
	Domain string
}

func (e *InvalidRecordDomainError) Error() string {
	return fmt.Sprintf("%s is set to the empty value but it is required", e.Domain)
}

func NewRecord(domain string, address string, ttl uint32) (*Record, error) {
	record := &Record{
		Domain:  domain,
		Address: address,
		TTL:     ttl,
	}

	err := ValidateRecord(record)
	if err != nil {
		return nil, fmt.Errorf("NewRecord: %w", err)
	}

	return record, nil
}

func ValidateRecord(record *Record) error {
	if !valid.IsDNSName(record.Domain) {
		return fmt.Errorf("ValidateRecord: %w", &InvalidRecordDomainError{
			Domain: record.Domain,
		})
	}

	if net.ParseIP(record.Address) == nil {
		return fmt.Errorf("ValidateRecord: %w", &InvalidRecordAddressError{
			Address: record.Address,
		})
	}

	return nil
}
