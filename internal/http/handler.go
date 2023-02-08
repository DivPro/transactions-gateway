package http

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	api "github.com/divpro/transactions-example/pkg/entity"
	"net/http"
)

type Handler struct {
	producer sarama.SyncProducer
}

func NewHandler(producer sarama.SyncProducer) Handler {
	return Handler{producer: producer}
}

func (h Handler) DepositCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request api.Deposit

	// TODO from auth
	request.UserID = "adeb0c95-a333-4049-9c27-beb34ec3cd32"
	request.UserID = "0a68a321-de8e-48ef-83b3-c256e087d310"

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _, err = h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "deposits",
		Value: sarama.ByteEncoder(data),
	})
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) TransactionCreate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request api.Transaction
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _, err = h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "transactions",
		Value: sarama.ByteEncoder(data),
	})
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) TransactionList(w http.ResponseWriter, req *http.Request) {
	// TODO
}

func (h Handler) UserList(w http.ResponseWriter, req *http.Request) {
	// TODO
}
