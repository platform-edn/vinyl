package dns_test

import (
	"errors"
	"net"
	"testing"

	miekg "github.com/miekg/dns"
	vinyl "github.com/platform-edn/vinyl/internal"
	"github.com/platform-edn/vinyl/internal/dns"
	"github.com/platform-edn/vinyl/internal/dns/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type QR struct {
	Question miekg.Question
	Record   vinyl.Record
}

type QRMap map[string]QR

func getRecordFunc(qrs QRMap) func(string) *vinyl.Record {
	return func(domain string) *vinyl.Record {
		q, exist := qrs[domain]
		if !exist {
			return nil
		}

		return &q.Record
	}
}

func inspectResponseFunc(rcode int, qrs QRMap, err error, assert *assert.Assertions) func(*miekg.Msg) error {
	return func(msg *miekg.Msg) error {
		if err != nil {
			return err
		}

		assert.Equal(rcode, msg.Rcode, "rcodes should be the same")

		if rcode == 2 {
			return nil
		}

		for _, rr := range msg.Answer {
			qr, exist := qrs[rr.Header().Name]

			assert.True(exist, "qr should exist")
			assert.Equal(qr.Record.Domain, rr.Header().Name, "domains should be the same")
			assert.Equal(qr.Record.TTL, rr.Header().Ttl, "ttl should be the same")
			assert.Equal(uint16(1), rr.Header().Rrtype, "rrtypes should be the same")

			arec := rr.(*miekg.A)
			assert.Equal(net.ParseIP(qr.Record.Address), arec.A, "addresses should be the same")
		}

		return nil
	}
}

func TestHandler_ServeDNS(t *testing.T) {
	tests := map[string]struct {
		QRs              QRMap
		ParseQuestion    bool
		ParseQuestionErr error
		Rcode            int
		Opcode           int
		WriteMsgErr      error
	}{
		"returns answers to questions correctly": {
			QRs: QRMap{
				"test.com": {
					Question: miekg.Question{
						Name:  "test.com",
						Qtype: 1,
					},
					Record: vinyl.Record{
						Domain:  "test.com",
						Address: "127.0.0.1",
						TTL:     3000,
					},
				},
			},
			ParseQuestion: true,
			Rcode:         0,
			Opcode:        0,
		},
		"errors if opcode isn't correct": {
			QRs: QRMap{
				"test.com": {
					Question: miekg.Question{
						Name:  "test.com",
						Qtype: 1,
					},
					Record: vinyl.Record{
						Domain:  "test.com",
						Address: "127.0.0.1",
						TTL:     3000,
					},
				},
			},
			ParseQuestion: false,
			Rcode:         2,
			Opcode:        1,
		},
		"errors if parseQuestion fails": {
			QRs: QRMap{
				"test!.com": {
					Question: miekg.Question{
						Name:  "test!.com",
						Qtype: 1,
					},
					Record: vinyl.Record{
						Domain:  "test@.com",
						Address: "127.0.0.1",
						TTL:     3000,
					},
				},
			},
			ParseQuestion: true,
			ParseQuestionErr: &vinyl.InvalidRecordDomainError{
				Domain: "test!.com",
			},
			Rcode:  2,
			Opcode: 0,
		},
		"writes error is writeMsg fails": {
			QRs: QRMap{
				"test.com": {
					Question: miekg.Question{
						Name:  "test.com",
						Qtype: 1,
					},
					Record: vinyl.Record{
						Domain:  "test.com",
						Address: "127.0.0.1",
						TTL:     3000,
					},
				},
			},
			ParseQuestion: true,
			WriteMsgErr:   errors.New("bad send"),
			Rcode:         0,
			Opcode:        0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			store := mocks.NewRecordStorer(t)

			if test.ParseQuestion {
				store.EXPECT().GetRecord(mock.AnythingOfType("string")).Call.Return(getRecordFunc(test.QRs), test.ParseQuestionErr)
			}

			rw := mocks.NewResponseWriter(t)

			rw.EXPECT().WriteMsg(mock.Anything).Call.Return(inspectResponseFunc(test.Rcode, test.QRs, test.WriteMsgErr, assert))

			handler := dns.NewRecordHandler(store)
			questions := []miekg.Question{}

			for _, qr := range test.QRs {
				questions = append(questions, qr.Question)
			}

			req := &miekg.Msg{
				MsgHdr: miekg.MsgHdr{
					Id:               123,
					RecursionDesired: false,
					CheckingDisabled: true,
					Opcode:           test.Opcode,
				},
				Question: questions,
			}

			handler.ServeDNS(rw, req)
		})
	}
}

func TestHandler_ServerErrorResponse(t *testing.T) {

	assert := assert.New(t)
	store := mocks.NewRecordStorer(t)
	handler := dns.NewRecordHandler(store)
	rw := mocks.NewResponseWriter(t)
	rw.EXPECT().WriteMsg(mock.Anything).Call.Return(inspectResponseFunc(0, nil, errors.New("bad"), assert))

	handler.ServerErrorResponse(rw, &miekg.Msg{}, errors.New("bad"))

}

func TestHandler_ParseQuestion(t *testing.T) {
	tests := map[string]struct {
		QRs             QRMap
		BadRecord       bool
		BadRecordType   uint16
		BadRecordDomain string
		GetRecord       bool
		GetRecordErr    error
		Err             error
	}{
		"parses A records correctly": {
			QRs: QRMap{
				"test.com": {
					Question: miekg.Question{
						Name:  "test.com",
						Qtype: 1,
					},
					Record: vinyl.Record{
						Domain:  "test.com",
						Address: "127.0.0.1",
						TTL:     3000,
					},
				},
			},
			GetRecord:    true,
			GetRecordErr: nil,
			Err:          nil,
		},
		"returns error if record isn't stored": {
			QRs:             QRMap{},
			BadRecord:       true,
			BadRecordDomain: "bad.com",
			BadRecordType:   1,
			GetRecord:       true,
			GetRecordErr:    errors.New("missing record"),
		},
		"returns error if record isn't type A": {
			QRs:             QRMap{},
			BadRecord:       true,
			BadRecordDomain: "bad.com",
			BadRecordType:   2,
			Err: &dns.UnsupportedRecordTypeError{
				Type: 2,
			},
		},
		"returns error if record can't be turned to A record": {
			QRs: QRMap{
				"test!.com": {
					Question: miekg.Question{
						Name:  "test!.com",
						Qtype: 1,
					},
					Record: vinyl.Record{
						Domain:  "test!.com",
						Address: "127.0.0.1",
						TTL:     3000,
					},
				},
			},
			GetRecord: true,
			Err: &vinyl.InvalidRecordDomainError{
				Domain: "test!.com",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			store := mocks.NewRecordStorer(t)

			if test.GetRecord {
				store.EXPECT().GetRecord(mock.AnythingOfType("string")).Call.Return(getRecordFunc(test.QRs), test.GetRecordErr)
			}

			handler := dns.NewRecordHandler(store)
			questions := []miekg.Question{}

			for _, qr := range test.QRs {
				questions = append(questions, qr.Question)
			}

			if test.BadRecord {
				questions = append(questions, miekg.Question{
					Name:  test.BadRecordDomain,
					Qtype: test.BadRecordType,
				})
			}

			rrs, err := handler.ParseQuestion(questions)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "error should be the same")
				return
			}
			if test.GetRecordErr != nil {
				assert.ErrorContains(err, test.GetRecordErr.Error(), "error should be the same")
				return
			}

			for _, rr := range rrs {
				qr, exist := test.QRs[rr.Header().Name]

				assert.True(exist, "qr should exist")
				assert.Equal(qr.Record.Domain, rr.Header().Name, "domains should be the same")
				assert.Equal(qr.Record.TTL, rr.Header().Ttl, "ttl should be the same")
				assert.Equal(uint16(1), rr.Header().Rrtype, "rrtypes should be the same")

				arec := rr.(*miekg.A)
				assert.Equal(net.ParseIP(qr.Record.Address), arec.A, "addresses should be the same")
			}
		})
	}
}
