package discovery

import (
	"context"
	"fmt"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/proto"
)

type RecordsServer struct {
	Store RecordStorer
	proto.UnimplementedRecordsServer
}

func NewRecordsServer(store RecordStorer) *RecordsServer {
	return &RecordsServer{
		Store: store,
	}
}

func (server *RecordsServer) CreateRecord(ctx context.Context, req *proto.CreateRecordRequest) (*proto.CreateRecordResponse, error) {
	record, err := server.Store.CreateRecord(req.Domain, req.Address, req.Ttl)
	if err != nil {
		return nil, fmt.Errorf("CreateRecord: %w", err)
	}

	resp := &proto.CreateRecordResponse{
		Record: &proto.Record{
			Domain:  record.Domain,
			Address: record.Address,
			Ttl:     record.TTL,
		},
	}

	return resp, nil
}

func (server *RecordsServer) RemoveRecord(ctx context.Context, req *proto.RemoveRecordRequest) (*proto.RemoveRecordResponse, error) {
	record, err := server.Store.RemoveRecord(req.Domain)
	if err != nil {
		return nil, fmt.Errorf("RemoveRecord: %w", err)
	}

	resp := &proto.RemoveRecordResponse{
		Record: &proto.Record{
			Domain:  record.Domain,
			Address: record.Address,
			Ttl:     record.TTL,
		},
	}

	return resp, nil
}

func (server *RecordsServer) GetRecord(ctx context.Context, req *proto.GetRecordRequest) (*proto.GetRecordResponse, error) {
	record, err := server.Store.GetRecord(req.Domain)
	if err != nil {
		return nil, fmt.Errorf("GetRecord: %w", err)
	}

	resp := &proto.GetRecordResponse{
		Record: &proto.Record{
			Domain:  record.Domain,
			Address: record.Address,
			Ttl:     record.TTL,
		},
	}

	return resp, nil
}

func (server *RecordsServer) ListRecords(sctx context.Context, req *proto.ListRecordsRequest) (*proto.ListRecordsResponse, error) {
	records, err := server.Store.ListRecords()
	if err != nil {
		return nil, fmt.Errorf("ListRecords: %w", err)
	}

	protoRecords := convertRecordsToProto(records...)

	resp := &proto.ListRecordsResponse{
		Records: protoRecords,
	}

	return resp, nil
}

func convertRecordsToProto(records ...vinyl.Record) []*proto.Record {
	protoRecords := []*proto.Record{}

	for _, record := range records {
		pr := &proto.Record{
			Domain:  record.Domain,
			Address: record.Address,
			Ttl:     record.TTL,
		}

		protoRecords = append(protoRecords, pr)
	}

	return protoRecords
}
