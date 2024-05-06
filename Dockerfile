FROM golang:1.21-alpine as builder
ARG TOKEN=""

RUN apk add --no-cache ca-certificates git

WORKDIR /currency-service
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install ./cmd/currency-service

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin /bin
USER nobody:nobody
ENTRYPOINT ["/bin/currency-service"]