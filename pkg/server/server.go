package server

import (
	"context"
	"encoding/json"
	"interview/configs"
	"interview/pkg/cache"
	"log"
	"net/http"
	"time"
)

//Server is a struct for sharing APIkey for methods
type Server struct {
	serv   *http.Server
	APIkey string
	Cache  *cache.Cache
}

var (
	//function that checks if you can do request right now
	allowRequest func() <-chan struct{}
)

//StartServer starts api server
func StartServer(ctx context.Context, config *configs.Config) error {
	server := new(Server)
	server.serv = &http.Server{
		Addr:    config.Addr,
		Handler: server.configureMux(),
	}
	server.APIkey = config.APIKey
	server.Cache = 	cache.NewCache()
	server.Cache.ReadFile()

	//go server.limitRequests()
	go func() {
		if err := server.serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(err.Error())
		}
	}()

	log.Printf("Starting server...")
	<-ctx.Done()
	log.Printf("Server stopped")

	server.Cache.WriteFile()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.serv.Shutdown(ctxShutDown); err != nil {
		log.Fatal("Server shutdown failed: ", err.Error())
	}

	log.Printf("Server exited properly")

	return nil
}

//configureMux sets routes
func (s *Server) configureMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/block/", s.handleCheckPath(http.HandlerFunc(s.handleSum())))
	return mux
}

//function for responding on request with json
func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

//limitRequests limiting requests amount on 3-rd party API
func (s *Server) limitRequests() {
	ticker := time.NewTicker(1 * time.Second)
	requestsLeft := make(chan struct{}, 1)
	allowRequest = s.allowRequest(requestsLeft)
	for {
		select {
		case <-ticker.C:
			if len(requestsLeft) < 1 {
				requestsLeft <- struct{}{}
			}
		}
	}
}

//allowRequest shows if request is allowed
func (s *Server) allowRequest(requestsLeft chan struct{}) func() <-chan struct{} {
	return func() <-chan struct{} {
		return requestsLeft
	}
}


