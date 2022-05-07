package dns_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/platform-edn/vinyl/internal/dns"
	"github.com/platform-edn/vinyl/internal/dns/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	assert := assert.New(t)
	handler := mocks.NewHandler(t)
	port := 53
	protocol := "udp"

	server := dns.NewServer(handler, port, protocol)

	assert.Contains(server.Addr, fmt.Sprint(port), "should contain assigned port")
	assert.Equal(server.Net, protocol, "should be the same protocol")
}

func TestStart(t *testing.T) {
	tests := map[string]struct {
		Err error
	}{
		"exits without error": {
			Err: nil,
		},
		"exits with error": {
			Err: errors.New("bad error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			server := mocks.NewServer(t)

			server.EXPECT().ListenAndServe().Return(test.Err)

			err := dns.Start(server)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "error should be the same")
				return
			}

			assert.NoError(err)
		})
	}
}
