ARG BUILDIMAGE=golang:1.17-alpine

FROM $BUILDIMAGE
RUN apk add --no-cache gcc musl-dev binutils-gold make git
RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.3