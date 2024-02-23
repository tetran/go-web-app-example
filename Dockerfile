# COntainer for building the app
FROM golang:1.21.1-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# -----------------------------------------------

# Container for deployment
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD ["./app"]

# -----------------------------------------------

# Hot reload container for development in local
FROM golang:1.21.1 as dev
WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.49.0
CMD ["air"]
