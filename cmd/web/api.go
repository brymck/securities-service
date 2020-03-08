package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	av "github.com/brymck/alpha-vantage-service/genproto/brymck/alpha_vantage/v1"
)

var alphaVantageApi av.AlphaVantageAPIClient

func getServiceAddress(serviceName string) string {
	return fmt.Sprintf("%s-4tt23pryoq-an.a.run.app:443", serviceName)
}

func init() {
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(getServiceAddress("alpha-vantage-service"), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	alphaVantageApi = av.NewAlphaVantageAPIClient(conn)
}
