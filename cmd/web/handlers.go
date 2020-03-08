package main

import (
	"context"
	"errors"
	"strconv"
	"time"

	av "github.com/brymck/alpha-vantage-service/genproto/brymck/alpha_vantage/v1"
	dt "github.com/brymck/genproto/brymck/dates/v1"
	sec "github.com/brymck/genproto/brymck/securities/v1"
	"github.com/brymck/helpers/dates"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brymck/securities-service/pkg/models"
)

func (app *application) GetSecurity(ctx context.Context, in *sec.GetSecurityRequest) (*sec.GetSecurityResponse, error) {
	response := &sec.GetSecurityResponse{}

	key := strconv.FormatUint(in.Id, 16)
	if err := app.getCache(key, response); err == nil {
		return response, nil
	}

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
		resp, err := alphaVantageApi.GetQuote(ctx, &av.GetQuoteRequest{Symbol: s.Symbol})
		if err != nil {
			return nil, err
		}
		s.Price = resp.Price
	}

	response = &sec.GetSecurityResponse{
		Security: &sec.Security{
			Id:     s.ID,
			Symbol: s.Symbol,
			Name:   s.Name,
			Price:  s.Price,
		},
	}
	_ = app.setCache(key, response)
	return response, nil
}

func (app *application) InsertSecurity(_ context.Context, in *sec.InsertSecurityRequest) (*sec.InsertSecurityResponse, error) {
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

	return &sec.InsertSecurityResponse{Id: id}, nil
}

func (app *application) GetPrices(_ context.Context, in *sec.GetPricesRequest) (*sec.GetPricesResponse, error) {
	startDate := dates.ProtoDateToTime(in.StartDate)
	endDate := dates.ProtoDateToTime(in.EndDate)
	log.Infof("requesting prices for %d between %s and %s", in.Id, dates.IsoFormat(startDate), dates.IsoFormat(endDate))
	records, err := app.prices.GetMany(&startDate, &endDate, in.Id, 1)
	if err != nil {
		return nil, err
	}
	log.Infof("retrieved %d records", len(records))
	var prices []*sec.Price
	for _, item := range records {
		year, month, day := item.Date.Date()
		date := dt.Date{Year: int32(year), Month: int32(month), Day: int32(day)}
		price := sec.Price{Date: &date, Price: item.Price}
		prices = append(prices, &price)
	}
	log.Infof("responding with %d price entries", len(prices))
	return &sec.GetPricesResponse{Prices: prices}, nil
}

func (app *application) UpdatePrices(ctx context.Context, in *sec.UpdatePricesRequest) (*sec.UpdatePricesResponse, error) {
	key := strconv.FormatUint(in.Id, 16)
	_ = app.cache.Delete(key)

	s, err := app.securities.Get(in.Id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			err := status.Error(codes.NotFound, err.Error())
			return nil, err
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp, err := alphaVantageApi.GetTimeSeries(ctx, &av.GetTimeSeriesRequest{Symbol: s.Symbol, Full: in.Full})
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

	return &sec.UpdatePricesResponse{Count: uint64(count)}, nil
}
