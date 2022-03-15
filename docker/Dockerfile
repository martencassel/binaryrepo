ARG BUILDIMAGE=golang:1.17-alpine

FROM $BUILDIMAGE
RUN apk add --no-cache gcc musl-dev binutils-gold make git bash
RUN go install github.com/go-delve/delve/cmd/dlv@v1.7.3 && \
    go install -v github.com/ramya-rao-a/go-outline@latest && \
    go install -v golang.org/x/tools/gopls@latest

COPY ./docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod u+x /docker-entrypoint.sh

ENTRYPOINT ["bash", "/docker-entrypoint.sh"]
CMD ["/src/build/binaryrepo", "run"]
