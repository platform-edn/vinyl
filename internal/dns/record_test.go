package dns_test

import (
	"net"
	"testing"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/dns"
	"github.com/stretchr/testify/assert"
)

func TestNewARecord(t *testing.T) {
	tests := map[string]struct {
		Record *vinyl.Record
		Err    error
	}{
		"creates a new A record": {
			Record: &vinyl.Record{
				Domain:  "test.com",
				Address: "127.0.0.1",
				TTL:     1000,
			},
			Err: nil,
		},
		"returns an error for bad domain": {
			Record: &vinyl.Record{
				Domain:  "test!.com",
				Address: "127.0.0.1",
				TTL:     1000,
			},
			Err: &vinyl.InvalidRecordDomainError{
				Domain: "test!.com",
			},
		},
		"returns an error for bad address": {
			Record: &vinyl.Record{
				Domain:  "test.com",
				Address: "1270.0.0.1",
				TTL:     1000,
			},
			Err: &vinyl.InvalidRecordAddressError{
				Address: "1270.0.0.1",
			},
		},
		"returns an error for bad ttl": {
			Record: &vinyl.Record{
				Domain:  "test.com",
				Address: "127.0.0.1",
				TTL:     0,
			},
			Err: &vinyl.InvalidRecordTTLError{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			arec, err := dns.NewARecord(test.Record)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
				return
			}

			assert.Equal(test.Record.Domain, arec.Hdr.Name)
			assert.Equal(net.ParseIP(test.Record.Address), arec.A)
			assert.Equal(test.Record.TTL, arec.Hdr.Ttl)
		})
	}
}
