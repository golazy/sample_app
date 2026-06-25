FROM golang:1.26-alpine AS build

WORKDIR /src

COPY . .
RUN go mod download
RUN go build -buildvcs=false -o /out/sample-app ./cmd/app

FROM alpine:3.22

RUN adduser -D -H -u 10001 app
USER app
WORKDIR /app

ENV ADDR=0.0.0.0:3000
EXPOSE 3000

COPY --from=build /out/sample-app /app/sample-app

CMD ["/app/sample-app"]
