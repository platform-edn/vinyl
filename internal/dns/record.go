package dns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	vinyl "github.com/platform-edn/vinyl/internal"
)

func NewARecord(record *vinyl.Record) (*dns.A, error) {
	err := vinyl.ValidateRecord(record)
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
