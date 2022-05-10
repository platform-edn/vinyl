package client

import (
	"context"
	"fmt"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/proto"
	"google.golang.org/grpc"
)

type Clienter interface {
	CreateRecord(ctx context.Context, in *proto.CreateRecordRequest, opts ...grpc.CallOption) (*proto.CreateRecordResponse, error)
	RemoveRecord(ctx context.Context, in *proto.RemoveRecordRequest, opts ...grpc.CallOption) (*proto.RemoveRecordResponse, error)
	GetRecord(ctx context.Context, in *proto.GetRecordRequest, opts ...grpc.CallOption) (*proto.GetRecordResponse, error)
	ListRecords(ctx context.Context, in *proto.ListRecordsRequest, opts ...grpc.CallOption) (*proto.ListRecordsResponse, error)
}

type RecordsClient struct {
	proto.RecordsClient
	Options []grpc.CallOption
}

func NewRecordsClient(protoClient proto.RecordsClient, options ...grpc.CallOption) *RecordsClient {
	return &RecordsClient{
		RecordsClient: protoClient,
		Options:       options,
	}
}

func (client *RecordsClient) Create(ctx context.Context, domain string, address string, ttl uint32) (*vinyl.Record, error) {
	// used to validate record before sending server side
	_, err := vinyl.NewRecord(domain, address, ttl)
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	resp, err := client.CreateRecord(
		ctx,
		&proto.CreateRecordRequest{
			Domain:  domain,
			Address: address,
			Ttl:     ttl,
		},
		client.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	record := &convertProtoToRecords(resp.Record)[0]

	return record, nil
}

func (client *RecordsClient) Remove(ctx context.Context, domain string) (*vinyl.Record, error) {
	resp, err := client.RemoveRecord(
		ctx,
		&proto.RemoveRecordRequest{
			Domain: domain,
		},
		client.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("Remove: %w", err)
	}

	record := &vinyl.Record{
		Domain:  resp.Record.Domain,
		Address: resp.Record.Address,
		TTL:     resp.Record.Ttl,
	}

	return record, nil
}

func (client *RecordsClient) Get(ctx context.Context, domain string) (*vinyl.Record, error) {
	resp, err := client.GetRecord(
		ctx,
		&proto.GetRecordRequest{
			Domain: domain,
		},
		client.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("Get: %w", err)
	}

	record := &vinyl.Record{
		Domain:  resp.Record.Domain,
		Address: resp.Record.Address,
		TTL:     resp.Record.Ttl,
	}

	return record, nil
}

func (client *RecordsClient) List(ctx context.Context) ([]vinyl.Record, error) {
	resp, err := client.ListRecords(
		ctx,
		&proto.ListRecordsRequest{},
		client.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("List: %w", err)
	}

	records := convertProtoToRecords(resp.Records...)

	return records, nil
}

func convertProtoToRecords(protoRecords ...*proto.Record) []vinyl.Record {
	records := []vinyl.Record{}

	for _, pr := range protoRecords {
		record := &vinyl.Record{
			Domain:  pr.Domain,
			Address: pr.Address,
			TTL:     pr.Ttl,
		}

		records = append(records, *record)
	}

	return records
}
