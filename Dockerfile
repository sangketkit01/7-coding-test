FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app cmd/api/*.go

FROM alpine:3.22
COPY --from=builder /app/app .

EXPOSE 8090
CMD [ "./app" ]
