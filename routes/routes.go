package routes

import (
	"context"
	pb "curr-conv-api/proto"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (server *HTTPServer) HandleGetCurrencies(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	currencyList, err := server.Client.GetCurrencyList(ctx, &pb.GetCurrencyListRequest{})
	if err != nil {
		SendMessage(res, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(res, http.StatusOK, currencyList)
}

func (server *HTTPServer) HandleConversion(res http.ResponseWriter, req *http.Request) {
	var data struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		SendMessage(res, fmt.Sprintf("Failed to retrieve data: %v.", err), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := server.Client.DoConvert(ctx, &pb.DoConvertRequest{
		From:   data.From,
		To:     data.To,
		Amount: data.Amount,
	})
	if err != nil {
		SendMessage(res, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(res, http.StatusOK, result)
}
