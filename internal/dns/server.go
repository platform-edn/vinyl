package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

type DNSServer struct {
	*dns.Server
}

func NewServer(handler dns.Handler, port int, protocol string) *DNSServer {
	srv := &dns.Server{Addr: fmt.Sprintf(":%v", port), Net: protocol}

	srv.Handler = handler
	server := &DNSServer{
		Server: srv,
	}

	return server
}

func (server *DNSServer) Start() error {
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("Start: %w", err)
	}

	return nil
}
