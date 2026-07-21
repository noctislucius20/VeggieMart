FROM golang:1.26.3-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o app main.go

EXPOSE 8000

CMD [ "./app" ]