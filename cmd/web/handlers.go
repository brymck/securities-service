package main

import (
	"context"
	"fmt"
	"time"

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

func getPrice(symbol string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	r, err := alphaVantageApi.GetQuote(ctx, &pb.GetQuoteRequest{Symbol: symbol})
	if err != nil {
		return 0.0, err
	}
	return r.Price, nil
}