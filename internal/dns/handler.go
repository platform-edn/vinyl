package dns

import (
	"fmt"
	"log"
	"net"

	"github.com/miekg/dns"
	vinyl "github.com/platform-edn/vinyl/internal"
)

type RecordStorer interface {
	GetRecord(string) (*vinyl.Record, error)
}

type ResponseWriter interface {
	// LocalAddr returns the net.Addr of the server
	LocalAddr() net.Addr
	// RemoteAddr returns the net.Addr of the client that sent the current request.
	RemoteAddr() net.Addr
	// WriteMsg writes a reply back to the client.
	WriteMsg(*dns.Msg) error
	// Write writes a raw buffer back to the client.
	Write([]byte) (int, error)
	// Close closes the connection.
	Close() error
	// TsigStatus returns the status of the Tsig.
	TsigStatus() error
	// TsigTimersOnly sets the tsig timers only boolean.
	TsigTimersOnly(bool)
	// Hijack lets the caller take over the connection.
	// After a call to Hijack(), the DNS package will not do anything with the connection.
	Hijack()
}

type RecordHandler struct {
	RecordStore RecordStorer
}

func NewRecordHandler(store RecordStorer) *RecordHandler {
	handler := &RecordHandler{
		RecordStore: store,
	}

	return handler
}

// ServeDNS implements interface for dns server handler
func (handler *RecordHandler) ServeDNS(w dns.ResponseWriter, request *dns.Msg) {
	response := NewResponse(request)
	log.Println(request)

	if request.Opcode != dns.OpcodeQuery {
		handler.ServerErrorResponse(w, response, fmt.Errorf("ServeDNS: %w", &UnsupportedOpCodeError{
			Opcode: request.Opcode,
		}))
		return
	}

	answers, err := handler.ParseQuestion(request.Question)
	if err != nil {
		handler.ServerErrorResponse(w, response, err)
		return
	}

	response.Answer = answers

	log.Println(response)

	err = w.WriteMsg(response)
	if err != nil {
		log.Println(err.Error())
	}
}

// ServerErrorResponse returns a server error and logs the error that caused it
func (handler *RecordHandler) ServerErrorResponse(w dns.ResponseWriter, response *dns.Msg, err error) {
	log.Println(err.Error())

	response.Rcode = 2

	err = w.WriteMsg(response)
	if err != nil {
		log.Println(err.Error())
	}
}

// ParseQuestion loops through questions and creates answers based on what is in the record store
func (handler *RecordHandler) ParseQuestion(questions []dns.Question) ([]dns.RR, error) {
	answers := []dns.RR{}

	for _, question := range questions {
		if question.Qtype != dns.TypeA {
			return nil, &UnsupportedRecordTypeError{
				Type: question.Qtype,
			}
		}

		record, err := handler.RecordStore.GetRecord(question.Name)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery: %w", err)
		}

		rr, err := NewARecord(record)
		if err != nil {
			return nil, fmt.Errorf("ParseQuery: %w", err)
		}

		answers = append(answers, rr)
	}

	return answers, nil
}

// NewResponse generates a response message based on a request message
func NewResponse(req *dns.Msg) *dns.Msg {
	resp := new(dns.Msg)
	resp.SetReply(req)
	resp.Compress = false

	return resp
}
