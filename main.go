package main

import (
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
)

type TokenConfig struct {
	User   string   `yaml:"user"`
	Groups []string `yaml:"groups"`
}

type Config struct {
	Tokens map[string]TokenConfig `yaml:"tokens"`
}

func (c Config) FindToken(t string) (*TokenConfig, error) {
	token, ok := c.Tokens[t]
	if !ok {
		return nil, fmt.Errorf("unable to find token")
	}

	return &token, nil
}

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}

	fs := flag.NewFlagSet("skaws", flag.ExitOnError)
	var (
		listenAddress = fs.String("listen-addr", ":8080", "listen address (also via LISTEN_ADDRESS)")
		configPath    = fs.String("config-path", "./example.yaml", "config path (also via CONFIG_PATH)")
	)
	err = ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatal("failed to parse config", zap.Error(err))
	}

	log.Info("loading config", zap.String("path", *configPath))
	c := Config{}
	configFileBytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal("failed to load config file", zap.Error(err))
	}

	err = yaml.Unmarshal(configFileBytes, &c)
	if err != nil {
		log.Fatal("failed to load config file", zap.Error(err))
	}

	h := &Handler{
		config: c,
		log:    log,
	}

	mux := http.NewServeMux()
	mux.Handle("/authenticate", h)

	log.Info("listening", zap.String("address", *listenAddress))
	err = http.ListenAndServe(*listenAddress, mux)
	if err != nil {
		log.Error("http listen err", zap.Error(err))
	}
}
