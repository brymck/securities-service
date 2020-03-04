package main

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brymck/securities-service/genproto"
	"github.com/brymck/securities-service/pkg/models"
)

func (app *application) GetSecurity(ctx context.Context, in *pb.GetSecurityRequest) (*pb.GetSecurityResponse, error) {
	s, err := app.securities.Get(in.Id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			err := status.Error(codes.NotFound, err.Error())
			return nil, err
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	if s.Price == 0.0 {
		log.Infof("retrieving missing price for %s", s.Symbol)
		resp, err := alphaVantageApi.GetQuote(ctx, &pb.GetQuoteRequest{Symbol: s.Symbol})
		if err != nil {
			return nil, err
		}
		s.Price = resp.Price
	}

	return &pb.GetSecurityResponse{
		Security: &pb.Security{
			Id:     s.ID,
			Symbol: s.Symbol,
			Name:   s.Name,
			Price:  s.Price,
		},
	}, nil
}

func (app *application) InsertSecurity(_ context.Context, in *pb.InsertSecurityRequest) (*pb.InsertSecurityResponse, error) {
	s := in.Security
	if s.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "symbol cannot be blank")
	}

	if s.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be blank")
	}

	id, err := app.securities.Insert(s.Symbol, s.Name)
	if err != nil {
		return nil, err
	}

	return &pb.InsertSecurityResponse{Id: id}, nil
}

func (app *application) UpdatePrices(ctx context.Context, in *pb.UpdatePricesRequest) (*pb.UpdatePricesResponse, error) {
	s, err := app.securities.Get(in.Id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			err := status.Error(codes.NotFound, err.Error())
			return nil, err
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp, err := alphaVantageApi.GetTimeSeries(ctx, &pb.GetTimeSeriesRequest{Symbol: s.Symbol, Full: in.Full})
	if err != nil {
		return nil, err
	}

	count := len(resp.TimeSeries)
	start := time.Now()
	log.Infof("inserting %d historical prices", count)
	for _, item := range resp.TimeSeries {
		date := time.Date(int(item.Date.Year), time.Month(item.Date.Month), int(item.Date.Day), 0, 0, 0, 0, time.UTC)
		err = app.prices.Insert(&date, s.ID, 1, item.Close)
		if err != nil {
			return nil, err
		}
	}
	end := time.Now()
	log.Infof("inserted %d historical prices in %d ms", count, end.Sub(start).Milliseconds())

	return &pb.UpdatePricesResponse{Count: int32(count)}, nil
}
