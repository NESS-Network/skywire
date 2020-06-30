package arclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/SkycoinProject/dmsg/cipher"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skywire-mainnet/internal/httpauth"
	"github.com/SkycoinProject/skywire-mainnet/pkg/snet/transport/tptypes"
)

var testPubKey, testSecKey = cipher.GenerateKeyPair()
var log = logging.MustGetLogger("arclient_test")

func TestClientAuth(t *testing.T) {
	wg := sync.WaitGroup{}

	headerCh := make(chan http.Header, 1)
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch url := r.URL.String(); url {
			case "/":
				defer wg.Done()
				headerCh <- r.Header

			case fmt.Sprintf("/security/nonces/%s", testPubKey):
				if _, err := fmt.Fprintf(w, `{"edge": "%s", "next_nonce": 1}`, testPubKey); err != nil {
					t.Errorf("Failed to write nonce response: %w", err)
				}

			default:
				t.Errorf("Don't know how to handle URL = '%s'", url)
			}
		},
	))
	defer srv.Close()

	apiClient, err := NewHTTP(srv.URL, testPubKey, testSecKey)
	require.NoError(t, err)
	c := apiClient.(*client)

	wg.Add(1)
	_, err = c.Get(context.TODO(), "/")
	require.NoError(t, err)

	header := <-headerCh
	assert.Equal(t, testPubKey.Hex(), header.Get("SW-Public"))
	assert.Equal(t, "1", header.Get("SW-Nonce"))
	assert.NotEmpty(t, header.Get("SW-Sig")) // TODO: check for the right key

	wg.Wait()
}

func TestBind(t *testing.T) {
	urlCh := make(chan string, 1)
	srv := httptest.NewServer(authHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlCh <- r.URL.String()
	})))
	defer srv.Close()

	c, err := NewHTTP(srv.URL, testPubKey, testSecKey)
	require.NoError(t, err)

	err = c.Bind(context.TODO(), tptypes.STCPR, "1234")
	require.NoError(t, err)

	assert.Equal(t, "/bind/stcpr", <-urlCh)
}

func authHandler(next http.Handler) http.Handler {
	m := http.NewServeMux()
	m.Handle("/security/nonces/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := json.NewEncoder(w).Encode(&httpauth.NextNonceResponse{Edge: testPubKey, NextNonce: 1}); err != nil {
				log.WithError(err).Error("Failed to encode nonce response")
			}
		},
	))
	m.Handle("/", next)
	return m
}
