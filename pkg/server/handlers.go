package server

import (
	"encoding/json"
	"fmt"
	"interview/tools"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

//handleSum counts sum of transactions for provided block
func (s *Server) handleSum() func(w http.ResponseWriter, r *http.Request) {
	type Value struct {
		Value string `json:"value"`
	}
	type Transactions struct {
		Values []Value `json:"transactions"`
	}
	type Result struct {
		Result *Transactions `json:"result"`
	}
	type Response struct {
		Transactions int     `json:"transactions"`
		Amount       float64 `json:"amount"`
	}

	const (
		WEI = 1000000000000000000 //convert to ether
	)
	return func(w http.ResponseWriter, r *http.Request) {
		block := strings.Split(r.URL.Path, "/")[3]
		blockNumber, err := strconv.ParseInt(block, 10, 64)
		if err != nil {
			log.Print(err)
			http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
			return
		}

		item, found := s.Cache.Get(blockNumber)
		if found {
			response := &Response{item.Transactions, item.Total}
			s.respond(w, r, 200, response)
			return
		}

		//<-allowRequest()
		path := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0x%x&boolean=true&apikey=%s", blockNumber, s.APIkey)
		resp, err := http.Get(path)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body) //for reusing json body
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		res := &Result{}
		if err := json.Unmarshal(bodyBytes, &res); err != nil {
			log.Print(err)
			s.errorJSONdecoding(w, r, bodyBytes)
			return
		}

		var sum float64

		for _, v := range res.Result.Values {
			num, err := strconv.ParseUint(v.Value[2:], 16, 64) //0x skipped
			if err != nil {                                   //can be too big number
				num, _, err := big.ParseFloat(v.Value[2:], 16, 0, 0)
				if err != nil {
					log.Print(err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
				add, _ := num.Float64()
				sum += add
				continue
			}
			sum += float64(num)
		}

		response := &Response{len(res.Result.Values), sum / WEI}
		s.respond(w, r, 200, response)
		s.Cache.Set(blockNumber, response.Transactions, response.Amount)
		//	w.Write([]byte(fmt.Sprintf("block: %s transactions: %d total: %.6f", block, len(res.Result.Values), sum/WEI)))

	}

}

func (s *Server) handleCheckPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tools.MatchPath([]byte(r.URL.Path)) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})
}

//errorJSONdecoding processes possible errors
func (s *Server) errorJSONdecoding(w http.ResponseWriter, r *http.Request, resp []byte) {
	type Err struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	someError := &Err{}
	if err := json.Unmarshal(resp, &someError); err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	s.respond(w, r, http.StatusOK, someError)
}
