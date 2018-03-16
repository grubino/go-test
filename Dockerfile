FROM golang:1.10.0-stretch

COPY main.go ./src/main.go

CMD ["go", "run", "./src/main.go"]
