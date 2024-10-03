FROM golang:1.23
WORKDIR /ddns
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./ddns ./cmd

CMD ["./ddns"]