# modules
FROM golang:1.15 as modules

ADD go.mod go.sum /m/
WORKDIR /m
RUN go mod download

# linter
FROM golang:1.15 as tester

ENV VERSION 1.27.0
ENV CHECKSUM 8d345e4e88520e21c113d81978e89ad77fc5b13bfdf20e5bca86b83fc4261272

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

COPY --from=modules /go/pkg /go/pkg

RUN golangci-lint run -v cmd/... domain/... internal/...

RUN go test -timeout=10s -v --race ./...

# Intermediate stage: Build the binary
FROM golang:1.15 as builder

COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 10001 myapp

RUN mkdir -p /retail
ADD . /retail
WORKDIR /retail

RUN GOOS=linux GOARCH=amd64 make build/api

# Final stage: Run the binary
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd
USER myapp

COPY --from=builder /retail/bin/retail /retail

CMD ["/retail"]
