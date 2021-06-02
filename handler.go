package main

import (
	"encoding/json"
	"go.uber.org/zap"
	authentication "k8s.io/api/authentication/v1beta1"
	"net/http"
)

type Handler struct {
	config Config
	log    *zap.Logger
}

func (h *Handler) handleErr(w http.ResponseWriter, status int, err error) {
	h.log.Error("writing http err response", zap.Error(err), zap.Int("status", status))
	w.WriteHeader(status)
	enc := json.NewEncoder(w)

	res := authentication.TokenReview{}
	res.APIVersion = "authentication.k8s.io/v1beta1"
	res.Kind = "TokenReview"
	res.Status = authentication.TokenReviewStatus{Authenticated: false, Error: err.Error()}

	err = enc.Encode(res)
	if err != nil {
		h.log.Error("an error occured writing the error response", zap.Error(err))
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var tr authentication.TokenReview
	err := decoder.Decode(&tr)
	if err != nil {
		h.handleErr(w, http.StatusBadRequest, err)
		return
	}
	h.log.Info("token authentication request received", zap.Any("body", tr))

	t, err := h.config.FindToken(tr.Spec.Token)
	if err != nil {
		h.handleErr(w, http.StatusUnauthorized, err)
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

	h.log.Info("authenticated successfully, writing response", zap.Any("res", res))

	err = enc.Encode(res)
	if err != nil {
		h.log.Error("failed to write success response", zap.Error(err))
	}
}
