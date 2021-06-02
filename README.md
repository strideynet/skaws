# skaws

Static Kubernetes Authentication Webhook Server

## Overview 

Skaws is an implementation of a Webhook token auth provider for Kubernetes as
described by https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication

This implementation is static (e.g it loads possible valid tokens from a yaml
file).

I beg that you do not use this in production. This tool is designed for
experimenting with Webhook Token Authentication.

## Handy bits

```shell
curl -X POST \
  'http://localhost:8080/authenticate' \
  -H 'Content-Type: application/json; charset=utf-8' \
  -d '{
  "apiVersion": "authentication.k8s.io/v1",
  "kind": "TokenReview",
  "spec": {
    "token": "valid-token"
  }
}'
```