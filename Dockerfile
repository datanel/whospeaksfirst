FROM golang:1.16 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -a -tags "-w -s" .

FROM scratch
WORKDIR /app
COPY public/ ./public
COPY --from=builder /src/whospeaksfirst ./
ENTRYPOINT ["/app/whospeaksfirst"]
