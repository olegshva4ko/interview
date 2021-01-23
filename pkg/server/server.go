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

/*
resp, err := http.Get("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0xafa025&boolean=true&apikey=YourApiKeyToken")
	if err != nil {
		fmt.Println(err)
		return
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err.Error())
	// }

	type Value struct {
		Value string `json:"value"`
	}
	type Transactions struct {
		Values []Value `json:"transactions"`
	}
	type Result struct {
		Result *Transactions `json:"result"`
	}
	//	res := &Result{}
	res1 := &Result{}
	json.NewDecoder(resp.Body).Decode(&res1)
	resp.Body.Close()
	//fmt.Printf("RES1\n%v\n", res1.Result.Values)
	//json.Unmarshal(body, &res)
	//fmt.Println(res)
	var (
		sum int64
	)
	for _, v := range res1.Result.Values {
		//fmt.Println(v.Value)
		num, _ := strconv.ParseInt(v.Value[2:], 16, 64)
		sum += num

	}
	var (
		dur1 time.Duration
		dur2 time.Duration
		r    float64
		buff bytes.Buffer
	)
	sum = 12
	for i := 0; i < 1000; i++ {
		start := time.Now()
		r = float64(sum) / 1000000000000000000
		//fmt.Println(result)
		//fmt.Println("First: ", time.Since(start))
		dur1 += time.Since(start)
	}

	fmt.Println("First: ", dur1/1000)
	for i := 0; i < 1000; i++ {
		start := time.Now()

		s := fmt.Sprintf("%d", sum)
		if len(s) >= 18 {
			buff.WriteString(s[:len(s)-18])
			buff.WriteString(".")
			buff.WriteString(s[len(s)-18:])
		} else {
			buff.WriteString("0.")
			for i := 0; i < 18-len(s); i++ {
				buff.WriteString("0")
			}
			buff.WriteString(s)
		}
		//fmt.Println(buff.String())
		dur2 += time.Since(start)
	}

	fmt.Println("Second: ", dur2/1000)

	os.Exit(0)
	fmt.Println(r)
*/
