FROM golang:1.19.4-alpine3.17 as build

WORKDIR /app
COPY ./* /app/
RUN go build -o ss2mqtt

FROM alpine:3.17.0 as runtime

WORKDIR /app
COPY --from=build /app/ss2mqtt /app/

CMD ["/app/ss2mqtt"]