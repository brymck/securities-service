package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/brymck/securities-service/genproto"
)

var alphaVantageApi pb.AlphaVantageAPIClient

func getServiceAddress(serviceName string) string {
	return fmt.Sprintf("%s-4tt23pryoq-an.a.run.app:443", serviceName)
}

func init() {
	creds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(getServiceAddress("alpha-vantage-service"), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	alphaVantageApi = pb.NewAlphaVantageAPIClient(conn)
}
