FROM golang:1.21-alpine
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
RUN go build -o /app/bin/server .
CMD ["/app/bin/server", "-http", ":8080"]
