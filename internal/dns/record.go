package dns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	vinyl "github.com/platform-edn/vinyl/internal"
)

func NewARecord(record *vinyl.Record) (*dns.A, error) {
	var err error
	switch {
	case record.Domain == "":
		err = &MissingFieldError{
			Field: "Domain",
		}
	case record.TTL == 0:
		err = &MissingFieldError{
			Field: "TTL",
		}
	case record.Address == "":
		err = &MissingFieldError{
			Field: "IP",
		}
	}
	if err != nil {
		return nil, fmt.Errorf("NewARecord: %w", err)
	}

	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   record.Domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    record.TTL,
		},
		A: net.ParseIP(record.Address),
	}

	return rr, nil
}
