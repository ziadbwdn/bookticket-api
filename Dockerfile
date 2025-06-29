FROM golang:1.24.2-alpine

WORKDIR /app

COPY ./root-app/go.mod ./root-app/go.sum ./
RUN go mod download

COPY ./root-app/ .

RUN go build -o ./out/app ./cmd/server

EXPOSE 8080

CMD [ "./out/app" ]