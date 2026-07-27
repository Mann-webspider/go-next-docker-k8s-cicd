FROM  golang:tip-alpine3.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 8080

CMD ["./main"]

FROM dhi.io/alpine-base:3.24 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/main /main

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/main"]