# Build site
FROM golang:1.14 as mustache

RUN go get github.com/cbroglie/mustache/...

FROM node:13 as site

COPY --from=mustache /go/bin/mustache /bin

ADD web /web
WORKDIR /web
RUN ./gen_static.sh
RUN yarn build

# linter
FROM golang:1.14 as tester

ENV VERSION 1.26.0
ENV CHECKSUM 59b0e49a4578fea574648a2fd5174ed61644c667ea1a1b54b8082fde15ef94fd

RUN echo "${CHECKSUM}  golangci-lint-${VERSION}-linux-amd64.tar.gz" > CHECKSUM

RUN curl -L "https://github.com/golangci/golangci-lint/releases/download/v${VERSION}/golangci-lint-${VERSION}-linux-amd64.tar.gz" \
    -o golangci-lint-${VERSION}-linux-amd64.tar.gz \
    && shasum -a 256 -c CHECKSUM \
    && tar xvzf golangci-lint-${VERSION}-linux-amd64.tar.gz \
    --strip-components=1 \
    -C ./bin \
    golangci-lint-${VERSION}-linux-amd64/golangci-lint

RUN rm -rf CHECKSUM "golangci-lint-${VERSION}-linux-amd64.tar.gz"

RUN mkdir -p /retail
ADD . /retail
WORKDIR /retail

RUN golangci-lint run --issues-exit-code=1 --deadline=600s ./...

RUN go test -timeout=600s -v --race ./...

# modules
FROM golang:1.14 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.14 as builder

COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 10001 myapp

RUN mkdir -p /retail
ADD . /retail
WORKDIR /retail

RUN GOOS=linux GOARCH=amd64 make build

# Final stage: Run the binary
FROM scratch

ENV PORT 8080
ENV DIAG_PORT 8181

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd
USER myapp

COPY --from=builder /retail/bin/retail /retail

COPY --from=site /web/dist/ /web/dist/

EXPOSE $PORT
EXPOSE $DIAG_PORT

CMD ["/retail"]
