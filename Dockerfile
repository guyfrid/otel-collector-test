FROM golang:1.18-alpine

WORKDIR /app

COPY ./splitbatchprocessor .

COPY ./example.yaml /etc/example.yaml

RUN go build -o main .

# RUN touch run.txt

CMD ["./main", "--config",  "example.yaml"]

