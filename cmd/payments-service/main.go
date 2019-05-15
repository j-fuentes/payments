package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/j-fuentes/payments/internal/restapi/controllers"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":3000", "Address and port where to listen")
	flag.Set("logtostderr", "true")
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	// Mount routes
	r.HandleFunc("/payments", controllers.GetPayments)

	glog.Infof("Listening on %s", addr)
	glog.Fatal(http.ListenAndServe(addr, r))
}
