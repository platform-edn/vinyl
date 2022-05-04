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
	return fmt.Sprintf("%s is not a valid ip address", e.Address)
}

type InvalidRecordDomainError struct {
	Domain string
}

func (e *InvalidRecordDomainError) Error() string {
	return fmt.Sprintf("%s is not a valid domain name", e.Domain)
}

type InvalidRecordTTLError struct {
	TTL uint32
}

func (e *InvalidRecordTTLError) Error() string {
	return fmt.Sprintf("%v is not a valid tll", e.TTL)
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

	if record.TTL == 0 {
		return fmt.Errorf("ValidateRecord: %w", &InvalidRecordTTLError{
			TTL: record.TTL,
		})
	}

	return nil
}
