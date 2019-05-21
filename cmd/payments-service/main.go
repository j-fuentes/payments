package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/j-fuentes/payments/internal/fixtures"
	"github.com/j-fuentes/payments/internal/restapi"
	"github.com/j-fuentes/payments/internal/store"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":3000", "Address and port where to listen")
	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	// Seed some sample data
	payments, err := fixtures.LoadPayments("payments.json")
	if err != nil {
		panic(err)
	}

	server := restapi.NewPaymentsServer(store.NewVolatilePaymentsStore(payments.Data))

	glog.Fatal(server.Serve(addr))
}
