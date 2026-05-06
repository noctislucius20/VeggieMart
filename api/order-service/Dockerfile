FROM golang:1.25.4-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o app main.go

EXPOSE 8082

CMD [ "./app" ]