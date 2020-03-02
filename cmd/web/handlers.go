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

type DatePrice struct {
	Date  *time.Time
	Price float64
}

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := alphaVantageApi.GetQuote(ctx, &pb.GetQuoteRequest{Symbol: symbol})
	if err != nil {
		return 0.0, err
	}
	return r.Price, nil
}

func getHistoricalPrices(symbol string) ([]*DatePrice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := alphaVantageApi.GetTimeSeries(ctx, &pb.GetTimeSeriesRequest{Symbol: symbol, Full: true})
	if err != nil {
		return nil, nil
	}
	ts := make([]*DatePrice, len(r.TimeSeries))
	for i, item := range r.TimeSeries {
		date := time.Date(int(item.Date.Year), time.Month(item.Date.Month), int(item.Date.Day), 0, 0, 0, 0, time.UTC)
		ts[i] = &DatePrice{Date: &date, Price: item.Close}
	}
	return ts, nil
}
