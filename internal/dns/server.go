package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

type DNSServer struct {
	*dns.Server
}

type Server interface {
	ListenAndServe() error
}

type Handler interface {
	ServeDNS(dns.ResponseWriter, *dns.Msg)
}

func NewServer(handler dns.Handler, port int, protocol string) *DNSServer {
	srv := &dns.Server{Addr: fmt.Sprintf(":%v", port), Net: protocol}

	srv.Handler = handler
	server := &DNSServer{
		Server: srv,
	}

	return server
}

func Start(server Server) error {
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("Start: %w", err)
	}

	return nil
}
