package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/j-fuentes/payments/internal/restapi"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":3000", "Address and port where to listen")
	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	// Sample data
	genDesc := "hello world"
	server := restapi.NewPaymentsServer(store.NewVolatilePaymentsStore(
		[]*models.Payment{
			&models.Payment{ID: 1, Description: &genDesc},
			&models.Payment{ID: 2, Description: &genDesc},
		},
	))

	glog.Fatal(server.Serve(addr))
}
