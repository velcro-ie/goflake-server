package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bstick12/goflake"
	"github.com/gorilla/mux"
)

var key = os.Getenv("KEYBASE")
var generator = goflake.GoFlakeInstanceUsingUnique(key)

func main() {
	startServer()
}

func startServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ids", Id)
	router.Queries("count", "{count:[0-9]+}")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Id(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()
	countVar := values["count"]
	var count = 1
	if len(countVar) == 1 {
		count, _ = strconv.Atoi(countVar[0])
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	ids := []string{}
	for i := 0; i < count; i++ {
		ids = append(ids, generator.GetBase64UUID())
	}
	json.NewEncoder(w).Encode(ids)

}
