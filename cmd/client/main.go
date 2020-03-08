package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	sec "github.com/brymck/genproto/brymck/securities/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type tokenAuth struct {
	token string
}

func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", t.token),
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

func getConnection(addr string) (*grpc.ClientConn, error) {
	pool, _ := x509.SystemCertPool()
	ce := credentials.NewClientTLSFromCert(pool, "")

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(ce),
		grpc.WithPerRPCCredentials(tokenAuth{token: os.Getenv("BRYMCK_ID_TOKEN")}),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	// Contact the server and print out its response.
	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Set up a connection to the server.
	conn, err := getConnection(os.Getenv("ADDR"))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := sec.NewSecuritiesAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.GetSecurity(ctx, &sec.GetSecurityRequest{Id: uint64(id)})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("%v", r)
}
