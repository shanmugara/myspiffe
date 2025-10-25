package goclient

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

func GetTlsClient(ctx context.Context) (*http.Client, error) {
	// Use the caller-provided context directly. Creating a child context and
	// deferring its cancel() here would cancel the context when this function
	// returns, which would make the X509 source unusable by the returned client.

	addr := os.Getenv("SPIFFE_ENDPOINT_SOCKET")
	if addr == "" {
		addr = "unix:///tmp/agent.sock"
		Logger.Infof("Using default socket endpoint: %s", addr)
	} else {
		Logger.Infof("Using socket endpoint from environment: %s", addr)
	}

	if !strings.HasPrefix(addr, "unix:") {
		addr = "unix://" + addr
	}
	mySvid, err := workloadapi.FetchX509SVID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch X509 SVID: %w", err)
	}
	myTD := mySvid.ID.TrustDomain()

	Logger.Infof("Using trust domain: %s", myTD)
	Logger.Infof("Using workload API endpoint: %s", addr)
	Logger.Infof("SVID: %s", mySvid.ID.URL())

	// Create a `workloadapi.X509Source`, it will connect to Workload API using provided socket path
	// If socket path is not defined using `workloadapi.SourceOption`, value from environment variable `SPIFFE_ENDPOINT_SOCKET` is used.
	source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(addr)))
	if err != nil {
		return nil, fmt.Errorf("unable to create X509Source: %w", err)
	}
	// Do NOT close the source here (no defer source.Close()). The returned client's
	// TLS config uses the source for SVIDs and bundle updates; closing it before the
	// client is finished would break mTLS. The caller/process should keep the source
	// alive for the lifetime of the client (or we could change the API to return the
	// source so the caller can close it when appropriate).

	// Allow connection to my trust domain member servers
	tlsConfig := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeMemberOf(myTD))
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return client, nil
}
