package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/j-fuentes/payments/internal/restapi/controllers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/payments", controllers.GetPayments)
	log.Fatal(http.ListenAndServe(":3000", r))
}
