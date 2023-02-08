package main

import (
	"flag"
	"log"
	coreHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/divpro/transactions-gateway/internal/config"
	"github.com/divpro/transactions-gateway/internal/http"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "Configuration file name")
	flag.Parse()
}

func main() {
	textHandler := slog.NewTextHandler(os.Stdout)
	logger := slog.New(textHandler)

	f, err := os.Open(configPath)
	if err != nil {
		logger.Error("open configuration file", err, configPath)
		return
	}
	var conf config.Config
	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		logger.Error("parse configuration file", err, configPath)
		return
	}

	saramaConf := sarama.NewConfig()
	saramaConf.ClientID = "transactions-gateway"
	saramaConf.Producer.Return.Successes = true

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	producer, err := sarama.NewSyncProducer(conf.Kafka.Brokers, saramaConf)
	if err != nil {
		logger.Error("parse configuration file", err, configPath)
		return
	}
	handler := http.NewHandler(producer)

	openapiMW, err := http.Middleware()
	if err != nil {
		logger.Error("prepare openapi validator", err, configPath)
		return
	}
	router := mux.NewRouter()
	router.Use(
		openapiMW,
	)
	router.PathPrefix("/api").Handler(http.Docs("/"))
	router.HandleFunc("/deposits", handler.DepositCreate).Methods(coreHttp.MethodPost)
	router.HandleFunc("/deposits", handler.DepositCreate).Methods(coreHttp.MethodPost)
	router.HandleFunc("/transactions", handler.TransactionCreate).Methods(coreHttp.MethodPost)
	router.HandleFunc("/transactions", handler.TransactionList).Methods(coreHttp.MethodGet)
	router.HandleFunc("/users", handler.UserList).Methods(coreHttp.MethodGet)

	go func() {
		logger.Info("starting http server on :8083")
		srv := &coreHttp.Server{
			Handler:      router,
			Addr:         "0.0.0.0:8083",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Fatal(srv.ListenAndServe())
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm
	logger.Info("terminating")
}
