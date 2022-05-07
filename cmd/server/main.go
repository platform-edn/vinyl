package main

import (
	"fmt"
	"log"
	"net"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/discovery"
	"github.com/platform-edn/vinyl/internal/dns"
	"github.com/platform-edn/vinyl/internal/proto"
	"github.com/platform-edn/vinyl/internal/store"
	"google.golang.org/grpc"
)

type RecordStorer interface {
	CreateRecord(string, string, uint32) (*vinyl.Record, error)
	RemoveRecord(string) (*vinyl.Record, error)
	ListRecords() ([]vinyl.Record, error)
	GetRecord(string) (*vinyl.Record, error)
}

func main() {
	store := store.NewMemory()

	StartDNSServer(store, 53, "udp")
	// go StartAPIServer(store, 9005)
}

func StartDNSServer(store RecordStorer, port int, protocol string) {
	handler := dns.NewRecordHandler(store)
	server := dns.NewServer(handler, port, protocol)

	err := dns.Start(server)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func StartAPIServer(store RecordStorer, port int, options ...grpc.ServerOption) {
	recordServer := discovery.NewRecordsServer(store)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(options...)
	proto.RegisterRecordsServer(grpcServer, recordServer)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
