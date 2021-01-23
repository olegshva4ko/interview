package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	possibleErrors := []struct {
		expected []byte //expected result
		path     string //path for request
		repeat   int    //count of request repetitions
	}{
		// { //The 5 calls per sec/IP rate limit is exceeded:
		// 	[]byte(`{"status":"0","message":"NOTOK","result":"Max rate limit reached"}`),
		// 	"https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0xafa01b&boolean=true&apikey=MR3R5T2B1UVQFGMX31IWTWD6VRXKNCKXY9",
		// 	10,
		// },
		{ //The 5 calls per sec/IP rate limit is exceeded:
			//this error is not mentioned in official API documentation
			[]byte(`{"status":"0","message":"NOTOK","result":"Max rate limit reached, please use API Key for higher rate limit"}`),
			"https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0xafa01b&boolean=true&apikey=YourApiKeyToken",
			2,
		},
		// { //An API request with a blank API Key or the default "YourApiKeyToken":
		// 	[]byte(`{"status":"1","message":"OK-Missing/Invalid API Key, rate limit of 1/5sec applied","result":"595623370144773018344492"}`),
		// 	"https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0xafa01b&boolean=true",
		// 	0,
		// },
		{ //An API request with Invalid API Key:
			[]byte(`{"status":"0","message":"NOTOK","result":"Invalid API Key"}`),
			"https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0xafa01b&boolean=true&apikey=123",
			0,
		},
		{ //API requests with Invalid API Key exceeding limit
			[]byte(`{"status":"0","message":"NOTOK","result":"Too many invalid api key attempts, please try again later"}`),
			"https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=0xafa01b&boolean=true&apikey=123",
			10,
		},
	}


	for _, possibleError := range possibleErrors {

		for i := 0; i < possibleError.repeat; i++ {
			go http.Get(possibleError.path)
		}
		resp, err := http.Get(possibleError.path)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body) //for reusing json body
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if len(bodyBytes) != len(possibleError.expected) {
			fmt.Println(string(bodyBytes), string(possibleError.expected))
			t.Error(err)
			t.FailNow()
		}
		for i := range bodyBytes {
			if bodyBytes[i] != possibleError.expected[i] {
				t.Error(err)
				t.FailNow()
			}
		}
		time.Sleep(3 * time.Second)
	}
}
