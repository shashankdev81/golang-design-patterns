package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Test is what we usually do"))
}
func main() {
	http.HandleFunc("/", LazyServer)
	http.ListenAndServe(":1111", nil)

	router := mux.NewRouter()
	router.HandleFunc("/test", TestEndpoint).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

}

func LazyServer(w http.ResponseWriter, r *http.Request) {
	headOrTails := rand.Intn(2)

	if headOrTails == 0 {
		time.Sleep(6 * time.Second)
		fmt.Fprintf(w, "Go! Slow server %v", headOrTails)
		fmt.Printf("Go! Slow server %v", headOrTails)
		return
	}

	fmt.Fprintf(w, "Go! Quick server %v", headOrTails)
	fmt.Printf("Go! Quick server %v", headOrTails)
	return
}
