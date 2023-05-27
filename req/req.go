package req

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Response struct {
	Result float64
	Error  string
}

func Fetch(pair string) (Response, error) {
	URL := fmt.Sprintf(
		"https://api.cryptowat.ch/markets/kraken/%s/price", pair)

	r, err := http.Get(URL)
	if err != nil {
		return Response{},
			errors.New(
				fmt.Sprintf("Error while fetching pair: %s: %s", pair, err))
	}

	defer r.Body.Close()

	// parse json
	var Model resModel

	d := json.NewDecoder(r.Body)
	err = d.Decode(&Model)
	if err != nil {
		return Response{}, errors.New(
			fmt.Sprintf("parsing failed for pair %s: %s", pair, err))
	}

	finalRes := Model.Result

	return Response{
		Result: finalRes.Price,
		Error:  finalRes.Error,
	}, nil
}
