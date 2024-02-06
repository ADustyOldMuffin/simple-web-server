FROM golang:1.21.6-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /simple-web-server

FROM scratch as release-stage
WORKDIR /
COPY --from=builder /simple-web-server /simple-web-server
ENTRYPOINT [ "/simple-web-server" ]
