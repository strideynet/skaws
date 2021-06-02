# skaws

Static Kubernetes Authentication Webhook Server

Skaws is an implementation of a Webhook token auth provider for Kubernetes as
described by https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication

This implementation is static (e.g it loads possible valid tokens from a yaml
file).

I beg that you do not use this in production. This tool is designed for
experimenting with Webhook Token Authentication.
