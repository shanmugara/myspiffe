package goclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func GetCert(server string) error {
	if Logger == nil {
		Logger = logrus.New()
	}
	ctx := context.Background()
	client, err := GetTlsClient(ctx)
	if err != nil {
		Logger.Errorf("Failed to create TLS client: %v", err)
		return err
	}

	// marshal SampleCSR to JSON
	payload, err := json.Marshal(SampleCSR)
	if err != nil {
		Logger.Errorf("Failed to marshal SampleCSR: %v", err)
		return fmt.Errorf("failed to marshal CSR: %w", err)
	}

	r, err := client.Post(server, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error connecting to %q: %w", server, err)
	}

	defer func() {
		if cerr := r.Body.Close(); cerr != nil {
			if Logger != nil {
				Logger.Errorf("failed to close response body: %v", cerr)
			}
		}
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read body: %w", err)
	}

	Logger.Infof("Got response from server: %s", string(body))
	return nil
}
