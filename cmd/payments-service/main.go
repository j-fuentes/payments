package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/j-fuentes/payments/internal/restapi"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":3000", "Address and port where to listen")
	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	glog.Fatal(restapi.Serve(addr))
}
