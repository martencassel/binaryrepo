# builder image
ARG BUILDIMAGE=golang:1.17-alpine

FROM $BUILDIMAGE as builder
RUN apk add --no-cache gcc musl-dev binutils-gold make git bash
RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.3 && \
    go install -v github.com/ramya-rao-a/go-outline@latest && \
    go install -v golang.org/x/tools/gopls@latest
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN make server

# generate clean, final image for end users
FROM alpine:3.11.3
RUN apk update && apk add --no-cache postgresql-client
COPY --from=builder /build/build/binaryrepo .
CMD until pg_isready --host=binaryrepo-postgres; do sleep 1; done \
    && /binaryrepo run