package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type TokenConfig struct {
	User   string
	Groups []string
}

type Config struct {
	Tokens map[string]TokenConfig
}

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}

	c := Config{}
	// TODO: Fetch config from YAML

	mux := http.NewServeMux()

	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {

	})

	log.Info("listening", zap.Int("port", 8080))
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Error("http listen err", zap.Error(err))
	}
}
