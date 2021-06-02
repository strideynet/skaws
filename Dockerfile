# syntax = docker/dockerfile:1.0-experimental
# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go get

# Copy the go source
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o skaws ./

# Create a final image
FROM alpine:3.13.5

WORKDIR /
RUN addgroup --gid 1000 -S skaws && adduser -S skaws -G skaws --uid 1000

COPY --from=builder /workspace/skaws .

USER skaws:skaws
ENTRYPOINT ["/skaws"]
