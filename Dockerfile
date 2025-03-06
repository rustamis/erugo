FROM golang:1.19
WORKDIR /app
COPY . .
RUN go mod download && go build -o erugo .
CMD ["./erugo"]
