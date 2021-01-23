package server_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"
)

func TestHandlerSum(t *testing.T) {
	type Value struct {
		Value string `json:"value"`
	}
	type Transactions struct {
		Values []Value `json:"transactions"`
	}
	type Result struct {
		Result *Transactions `json:"result"`
	}
	const (
		WEI = 1000000000000000000 //convert to ether
	)

	bodies := []struct {
		body         string
		total        int
		transactions float64
	}{
		{
			`{"result": {"transactions" :[{"value": "0xFFFFFFFFFFFFFFFFF"}]}}`,
			1,
			float64(295147905179352825855) / float64(WEI),
		},
		{
			`{"result": {"transactions" :[{"value": "0xFFFFFFFFFFFFFFFFF"}, {"value": "0xFFFFFFFFFFFFFFFFF"}]}}`,
			2,
			float64(295147905179352825855 * 2) / float64(WEI),
		},
	}
	for i := range bodies {
		res := &Result{}
		if err := json.Unmarshal([]byte(bodies[i].body), &res); err != nil {
			t.Errorf(err.Error())
		}

		var sum float64

		for _, v := range res.Result.Values {
			num, err := strconv.ParseUint(v.Value[2:], 16, 64)
			if err != nil {
				num, _, err := big.ParseFloat(v.Value[2:], 16, 0, 0)
				if err != nil {
					t.Errorf(err.Error())

				}
				add, _ := num.Float64()
				sum += add
				continue
			}
			sum += float64(num)
		}

		response := fmt.Sprintf("Total: %d\nTransactions: %.5f\n", len(res.Result.Values), sum/WEI)
		expect := fmt.Sprintf("Total: %d\nTransactions: %.5f", bodies[i].total, bodies[i].transactions)
		if response != expect {
			fmt.Println(response, expect)
			t.Error()
		}
	}
}
