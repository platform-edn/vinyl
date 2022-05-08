package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/discovery"
	"github.com/platform-edn/vinyl/internal/dns"
	"github.com/platform-edn/vinyl/internal/proto"
	"github.com/platform-edn/vinyl/internal/store"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type RecordStorer interface {
	CreateRecord(string, string, uint32) (*vinyl.Record, error)
	RemoveRecord(string) (*vinyl.Record, error)
	ListRecords() ([]vinyl.Record, error)
	GetRecord(string) (*vinyl.Record, error)
}

func main() {
	// setup os signal trigger for shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// setup error group for shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	// business logic
	store := store.NewMemory()

	// generate grpc services
	recordService := discovery.NewRecordsServer(store)

	// generate servers
	grpcServer := grpc.NewServer()
	proto.RegisterRecordsServer(grpcServer, recordService)

	serveDNSFunc, dnsServer := ServeDNS(store, 53, "udp")
	serveGRPCFunc := ServeGRPC(grpcServer, 8080)

	// start servers
	errGroup.Go(serveDNSFunc)
	errGroup.Go(serveGRPCFunc)

	// wait for shutdown signals
	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}

	// will trigger errGroup to shutdown if os signal is what caused shutdown
	cancel()

	log.Println("attempting to shutdown servers...")

	// shutdown servers
	err := dnsServer.Shutdown()
	if err != nil {
		log.Println(err)
	}
	grpcServer.GracefulStop()

	// wait for servers to shutdown
	err = errGroup.Wait()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	log.Println("Good bye!")
}

func ServeDNS(store RecordStorer, port int, protocol string) (func() error, *dns.DNSServer) {
	handler := dns.NewRecordHandler(store)
	server := dns.NewServer(handler, port, protocol)

	serverFunc := func() error {
		log.Println("starting dns server...")
		err := dns.Start(server)
		if err != nil {
			return err
		}

		return nil
	}

	return serverFunc, server
}

func ServeGRPC(server *grpc.Server, port int) func() error {
	serverFunc := func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
		if err != nil {
			return err
		}

		log.Println("starting grpc server...")
		err = server.Serve(lis)
		if err != nil {
			return err
		}

		return nil
	}

	return serverFunc
}
