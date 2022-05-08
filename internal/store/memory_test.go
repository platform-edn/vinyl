package store_test

import (
	"testing"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestMemory_GetRecord(t *testing.T) {
	domain := "test.com"
	address := "127.0.0.1"
	var ttl uint32 = 60
	record := vinyl.Record{
		Domain:  domain,
		Address: address,
		TTL:     ttl,
	}

	tests := map[string]struct {
		Records []vinyl.Record
		Err     error
	}{
		"returns record correctly": {
			Records: []vinyl.Record{record},
			Err:     nil,
		},
		"returns MissingRecordError when store is empty": {
			Records: []vinyl.Record{},
			Err: &store.MissingRecordError{
				Domain: domain,
			},
		},
		"returns MissingRecordError with exisitng records": {
			Records: []vinyl.Record{
				{
					Domain:  "bad.com",
					Address: "127.0.0.1",
					TTL:     10,
				},
			},
			Err: &store.MissingRecordError{
				Domain: domain,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			mem := store.NewMemory(test.Records...)

			rec, err := mem.GetRecord(domain)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
				return
			}

			assert.NoError(err)
			assert.Equal(record, *rec)
		})
	}
}

func TestMemory_RemoveRecord(t *testing.T) {
	domain := "test.com"
	address := "127.0.0.1"
	var ttl uint32 = 60
	record := vinyl.Record{
		Domain:  domain,
		Address: address,
		TTL:     ttl,
	}

	tests := map[string]struct {
		Records []vinyl.Record
		Err     error
	}{
		"removes record correctly": {
			Records: []vinyl.Record{record},
			Err:     nil,
		},
		"returns MissingRecordError": {
			Records: []vinyl.Record{},
			Err: &store.MissingRecordError{
				Domain: domain,
			},
		},
		"returns MissingRecordError with exisitng records": {
			Records: []vinyl.Record{
				{
					Domain:  "bad.com",
					Address: "127.0.0.1",
					TTL:     10,
				},
			},
			Err: &store.MissingRecordError{
				Domain: domain,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			mem := store.NewMemory(test.Records...)

			rec, err := mem.RemoveRecord(domain)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
				return
			}

			assert.NoError(err)
			assert.Equal(record, *rec)
			assert.Len(mem.Records, len(test.Records)-1)
		})
	}
}

func TestMemory_ListRecords(t *testing.T) {
	domain1 := "test1.com"
	domain2 := "test2.com"
	address := "127.0.0.1"
	var ttl uint32 = 60
	record1 := vinyl.Record{
		Domain:  domain1,
		Address: address,
		TTL:     ttl,
	}
	record2 := vinyl.Record{
		Domain:  domain2,
		Address: address,
		TTL:     ttl,
	}

	tests := map[string]struct {
		Records []vinyl.Record
		Err     error
	}{
		"returns correct amount of records": {
			Records: []vinyl.Record{record1, record2},
			Err:     nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			mem := store.NewMemory(test.Records...)

			records, err := mem.ListRecords()
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
				return
			}

			assert.NoError(err)
			assert.Equal(test.Records, records)

		})
	}
}

func TestMemory_CreateRecord(t *testing.T) {
	tests := map[string]struct {
		Records []vinyl.Record
		Domain  string
		Address string
		TTL     uint32
		Err     error
	}{
		"creates record correctly": {
			Records: []vinyl.Record{},
			Domain:  "test.com",
			Address: "127.0.0.1",
			TTL:     100,
			Err:     nil,
		},
		"returns error with bad domain": {
			Records: []vinyl.Record{},
			Domain:  "test!!!!.com",
			Address: "127.0.0.1",
			TTL:     100,
			Err: &vinyl.InvalidRecordDomainError{
				Domain: "test!!!!.com",
			},
		},
		"returns error with bad address": {
			Records: []vinyl.Record{},
			Domain:  "test.com",
			Address: "127.0.0.10000",
			TTL:     100,
			Err: &vinyl.InvalidRecordAddressError{
				Address: "127.0.0.10000",
			},
		},
		"returns error with bad ttl": {
			Records: []vinyl.Record{},
			Domain:  "test.com",
			Address: "127.0.0.100",
			TTL:     0,
			Err:     &vinyl.InvalidRecordTTLError{},
		},
		"returns error about existing record": {
			Records: []vinyl.Record{
				{
					Domain:  "test.com",
					Address: "127.0.0.100",
					TTL:     0,
				},
			},
			Domain:  "test.com",
			Address: "127.0.0.100",
			TTL:     0,
			Err: &store.ExistingRecordError{
				Domain: "test.com",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			mem := store.NewMemory(test.Records...)

			rec, err := mem.CreateRecord(test.Domain, test.Address, test.TTL)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error())
				return
			}

			assert.NoError(err)
			assert.Equal(test.Domain, rec.Domain)
			assert.Equal(test.Address, rec.Address)
			assert.Equal(test.TTL, rec.TTL)
		})
	}
}
