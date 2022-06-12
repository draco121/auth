FROM golang:1.18 as build

WORKDIR /auth-service

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine

WORKDIR /app

COPY --from=build /auth-service ./

EXPOSE 8080

CMD [ "./app" ]
