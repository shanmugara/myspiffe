package main

import (
	"github.com/sirupsen/logrus"
	"myspiffe/goclient"
	"time"
)

const (
	ServerUrl = "https://omegaspire01.omegaworld.net:8082/v1/cert/request"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.Infof("Starting client")

	for {
		err := goclient.GetCert(ServerUrl)
		if err != nil {
			logger.Errorf("GetCert failed: %v", err)
		}
		time.Sleep(10 * time.Second)
	}
}
