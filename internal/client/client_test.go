package client_test

import (
	"errors"
	"testing"
	"time"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/client"
	"github.com/platform-edn/vinyl/internal/client/mocks"
	"github.com/platform-edn/vinyl/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestRecordsClient_Create(t *testing.T) {
	tests := map[string]struct {
		Domain          string
		Address         string
		Ttl             uint32
		CreateRecord    bool
		CreateRecordErr error
		Err             error
	}{
		"successfully create record": {
			Domain:          "test.com",
			Address:         "127.0.0.1",
			Ttl:             3000,
			CreateRecord:    true,
			CreateRecordErr: nil,
			Err:             nil,
		},
		"returns error from server side": {
			Domain:          "test.com",
			Address:         "127.0.0.1",
			Ttl:             3000,
			CreateRecord:    true,
			CreateRecordErr: errors.New("server side error"),
			Err:             nil,
		},
		"returns error from bad record data": {
			Domain:          "test!.com",
			Address:         "127.0.0.1",
			Ttl:             3000,
			CreateRecord:    false,
			CreateRecordErr: nil,
			Err: &vinyl.InvalidRecordDomainError{
				Domain: "test!.com",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			clienter := mocks.NewClienter(t)

			if test.CreateRecord {
				clienter.EXPECT().CreateRecord(
					mock.Anything,
					mock.Anything,
				).Return(
					&proto.CreateRecordResponse{
						Record: &proto.Record{
							Domain:  test.Domain,
							Address: test.Address,
							Ttl:     test.Ttl,
						},
					},
					test.CreateRecordErr,
				)
			}

			client := client.NewRecordsClient(clienter)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			record, err := client.Create(ctx, test.Domain, test.Address, test.Ttl)
			if test.CreateRecordErr != nil {
				assert.ErrorContains(err, test.CreateRecordErr.Error(), "CreateRecordErr should be the same")
				return
			}
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "Err should be the same")
				return
			}

			assert.Equal(test.Domain, record.Domain)
			assert.Equal(test.Address, record.Address)
			assert.Equal(test.Ttl, record.TTL)
		})
	}
}

func TestRecordsClient_Remove(t *testing.T) {
	tests := map[string]struct {
		Domain  string
		Address string
		Ttl     uint32
		Err     error
	}{
		"successfully removes record": {
			Domain:  "test.com",
			Address: "127.0.0.1",
			Ttl:     3000,
			Err:     nil,
		},
		"returns error from server side": {
			Domain:  "test.com",
			Address: "127.0.0.1",
			Ttl:     3000,
			Err:     errors.New("server side error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			clienter := mocks.NewClienter(t)

			clienter.EXPECT().RemoveRecord(
				mock.Anything,
				mock.Anything,
			).Return(
				&proto.RemoveRecordResponse{
					Record: &proto.Record{
						Domain:  test.Domain,
						Address: test.Address,
						Ttl:     test.Ttl,
					},
				},
				test.Err,
			)

			client := client.NewRecordsClient(clienter)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			record, err := client.Remove(ctx, test.Domain)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "Err should be the same")
				return
			}

			assert.Equal(test.Domain, record.Domain)
			assert.Equal(test.Address, record.Address)
			assert.Equal(test.Ttl, record.TTL)
		})
	}
}

func TestRecordsClient_Get(t *testing.T) {
	tests := map[string]struct {
		Domain  string
		Address string
		Ttl     uint32
		Err     error
	}{
		"successfully gets record": {
			Domain:  "test.com",
			Address: "127.0.0.1",
			Ttl:     3000,
			Err:     nil,
		},
		"returns error from server side": {
			Domain:  "test.com",
			Address: "127.0.0.1",
			Ttl:     3000,
			Err:     errors.New("server side error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			clienter := mocks.NewClienter(t)

			clienter.EXPECT().GetRecord(
				mock.Anything,
				mock.Anything,
			).Return(
				&proto.GetRecordResponse{
					Record: &proto.Record{
						Domain:  test.Domain,
						Address: test.Address,
						Ttl:     test.Ttl,
					},
				},
				test.Err,
			)

			client := client.NewRecordsClient(clienter)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			record, err := client.Get(ctx, test.Domain)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "Err should be the same")
				return
			}

			assert.Equal(test.Domain, record.Domain)
			assert.Equal(test.Address, record.Address)
			assert.Equal(test.Ttl, record.TTL)
		})
	}
}

func TestRecordsClient_List(t *testing.T) {
	tests := map[string]struct {
		Domain  string
		Address string
		Ttl     uint32
		Err     error
	}{
		"successfully create record": {
			Domain:  "test.com",
			Address: "127.0.0.1",
			Ttl:     3000,
			Err:     nil,
		},
		"returns error from server side": {
			Domain:  "test.com",
			Address: "127.0.0.1",
			Ttl:     3000,
			Err:     errors.New("server side error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			clienter := mocks.NewClienter(t)

			clienter.EXPECT().ListRecords(
				mock.Anything,
				mock.Anything,
			).Return(
				&proto.ListRecordsResponse{
					Records: []*proto.Record{
						{
							Domain:  test.Domain,
							Address: test.Address,
							Ttl:     test.Ttl,
						},
					},
				},
				test.Err,
			)

			client := client.NewRecordsClient(clienter)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			records, err := client.List(ctx)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "Err should be the same")
				return
			}

			assert.Equal(test.Domain, records[0].Domain)
			assert.Equal(test.Address, records[0].Address)
			assert.Equal(test.Ttl, records[0].TTL)
		})
	}
}
