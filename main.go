package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"go.uber.org/zap"
	authentication "k8s.io/api/authentication/v1beta1"
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

func handleErr(l *zap.Logger, w http.ResponseWriter, status int, err error) {
	l.Error("writing http err response", zap.Error(err), zap.Int("status", status))
	w.WriteHeader(status)
	enc := json.NewEncoder(w)

	res := authentication.TokenReview{}
	res.APIVersion = "authentication.k8s.io/v1beta1"
	res.Kind = "TokenReview"
	res.Status = authentication.TokenReviewStatus{Authenticated: false, Error: err.Error()}

	err = enc.Encode(res)
	if err != nil {
		l.Error("an error occured writing the error response", zap.Error(err))
	}
}

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}

	fs := flag.NewFlagSet("skaws", flag.ExitOnError)
	var (
		listenAddress = fs.String("listen-addr", ":8080", "listen address (also via LISTEN_ADDRESS)")
		configPath    = fs.String("config-path", "/tokens.yaml", "config path (also via CONFIG_PATH)")
	)
	err = ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatal("failed to parse config", zap.Error(err))
	}

	c := Config{}
	// TODO: Fetch config from YAML

	mux := http.NewServeMux()

	mux.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var tr authentication.TokenReview
		err := decoder.Decode(&tr)
		if err != nil {
			handleErr(log, w, http.StatusBadRequest, err)
			return
		}

		t, err := c.FindToken(tr.Spec.Token)
		if err != nil {
			handleErr(log, w, http.StatusUnauthorized, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)

		res := authentication.TokenReview{}
		res.APIVersion = "authentication.k8s.io/v1beta1"
		res.Kind = "TokenReview"
		res.Status = authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				Username: t.User,
				Groups:   t.Groups,
			},
		}

		err = enc.Encode(res)
		if err != nil {
			log.Error("failed to write success response", zap.Error(err))
		}
	})

	log.Info("listening", zap.String("address", *listenAddress))
	err = http.ListenAndServe(*listenAddress, mux)
	if err != nil {
		log.Error("http listen err", zap.Error(err))
	}
}
