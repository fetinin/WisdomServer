FROM golang:1.18-stretch as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 go build -o client -ldflags "-s -w" ./cmd/client/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/client /app/
ENTRYPOINT ["/app/client"]
