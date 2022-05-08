package discovery_test

import (
	"context"
	"errors"
	"testing"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/discovery"
	"github.com/platform-edn/vinyl/internal/discovery/mocks"
	"github.com/platform-edn/vinyl/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecordsServer_CreateRecord(t *testing.T) {
	tests := map[string]struct {
		Err error
	}{
		"should successfully return a domain": {
			Err: nil,
		},
		"should successfully return an error": {
			Err: errors.New("bad error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			store := mocks.NewRecordStorer(t)
			record := &vinyl.Record{
				Domain:  "test.com",
				Address: "127.0.0.1",
				TTL:     3000,
			}

			store.EXPECT().CreateRecord(
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("uint32"),
			).Return(
				record,
				test.Err,
			)

			server := discovery.NewRecordsServer(store)

			resp, err := server.CreateRecord(context.Background(), &proto.CreateRecordRequest{
				Domain:  record.Domain,
				Address: record.Address,
				Ttl:     record.TTL,
			})
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "error should be the same")
				return
			}

			assert.NoError(err, "should not have returned error")
			assert.Equal(resp.Record.Domain, record.Domain, "domains should be the same")
			assert.Equal(resp.Record.Address, record.Address, "addresses should be the same")
			assert.Equal(resp.Record.Ttl, record.TTL, "ttls should be the same")
		})
	}
}

func TestRecordsServer_RemoveRecord(t *testing.T) {
	tests := map[string]struct {
		Err error
	}{
		"should successfully return a domain": {
			Err: nil,
		},
		"should successfully return an error": {
			Err: errors.New("bad error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			store := mocks.NewRecordStorer(t)
			record := &vinyl.Record{
				Domain:  "test.com",
				Address: "127.0.0.1",
				TTL:     3000,
			}

			store.EXPECT().RemoveRecord(
				mock.AnythingOfType("string"),
			).Return(
				record,
				test.Err,
			)

			server := discovery.NewRecordsServer(store)

			resp, err := server.RemoveRecord(context.Background(), &proto.RemoveRecordRequest{
				Domain: record.Domain,
			})
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "error should be the same")
				return
			}

			assert.NoError(err, "should not have returned error")
			assert.Equal(resp.Record.Domain, record.Domain, "domains should be the same")
			assert.Equal(resp.Record.Address, record.Address, "addresses should be the same")
			assert.Equal(resp.Record.Ttl, record.TTL, "ttls should be the same")
		})
	}
}

func TestRecordsServer_GetRecord(t *testing.T) {
	tests := map[string]struct {
		Err error
	}{
		"should successfully return a domain": {
			Err: nil,
		},
		"should successfully return an error": {
			Err: errors.New("bad error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			store := mocks.NewRecordStorer(t)
			record := &vinyl.Record{
				Domain:  "test.com",
				Address: "127.0.0.1",
				TTL:     3000,
			}

			store.EXPECT().GetRecord(
				mock.AnythingOfType("string"),
			).Return(
				record,
				test.Err,
			)

			server := discovery.NewRecordsServer(store)

			resp, err := server.GetRecord(context.Background(), &proto.GetRecordRequest{
				Domain: record.Domain,
			})
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "error should be the same")
				return
			}

			assert.NoError(err, "should not have returned error")
			assert.Equal(resp.Record.Domain, record.Domain, "domains should be the same")
			assert.Equal(resp.Record.Address, record.Address, "addresses should be the same")
			assert.Equal(resp.Record.Ttl, record.TTL, "ttls should be the same")
		})
	}
}

func TestRecordsServer_ListRecords(t *testing.T) {
	tests := map[string]struct {
		Err error
	}{
		"should successfully return a domain": {
			Err: nil,
		},
		"should successfully return an error": {
			Err: errors.New("bad error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			store := mocks.NewRecordStorer(t)
			record := vinyl.Record{
				Domain:  "test.com",
				Address: "127.0.0.1",
				TTL:     3000,
			}

			store.EXPECT().ListRecords().Return(
				[]vinyl.Record{record},
				test.Err,
			)

			server := discovery.NewRecordsServer(store)

			resp, err := server.ListRecords(context.Background(), &proto.ListRecordsRequest{})
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "error should be the same")
				return
			}

			assert.NoError(err, "should not have returned error")
			assert.Equal(resp.Records[0].Domain, record.Domain, "domains should be the same")
			assert.Equal(resp.Records[0].Address, record.Address, "addresses should be the same")
			assert.Equal(resp.Records[0].Ttl, record.TTL, "ttls should be the same")
		})
	}
}
