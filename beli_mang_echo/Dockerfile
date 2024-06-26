FROM golang:1.22.3-alpine AS builder

WORKDIR /app

COPY . . 

ENV GOARCH=amd64
ENV GOOS=linux 

RUN go build -o main ./main.go

FROM alpine
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

# ENV DB_HOST=
# ENV DB_PORT=
# ENV DB_USER=
# ENV DB_PASSWORD=
# ENV DB_NAME=
# ENV DB_PARAMS=
# ENV JWT_SECRET=
# ENV BCRYPT_SALT=
# ENV APP_PORT=

CMD ["/app/main"]