package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"bitbucket.org/retirementplanio/go-simulation/simulation"
)

///////////////
// Utilities //
///////////////

type response map[string]interface{}

func (r response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

//////////
// Main //
//////////

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")

	if port != "" {
		flag.Set("bind", ":"+port)
	}

	goji.Get("/", root)
	goji.Get("/health", health)

	authenticated := web.New()
	authenticated.Use(secured)
	goji.Handle("/simulation", authenticated)
	authenticated.Post("/simulation", simulateHandler)

	log.Println("Booting retirement simulation server...")
	goji.Serve()
}

////////////////
// Middleware //
////////////////

func secured(c *web.C, h http.Handler) http.Handler {
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		panic("AUTH_TOKEN environment variable required.")
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		providedToken := r.Header.Get("Authorization")

		if (providedToken != authToken) || (providedToken == "") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, response{"success": false, "message": "You must be authorized to perform that action."})
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

////////////////////
// Route Handlers //
////////////////////

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Go Simulation API")
	return
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
	return
}

func simulateHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	apiResponse := simulation.ValidateAndHandleJsonInput(r)
	end := time.Since(start)

	log.Printf("Processing request from %s in %vs", r.RemoteAddr, end)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiResponse.StatusCode)
	fmt.Fprint(w, response(apiResponse.Response))
	return
}
