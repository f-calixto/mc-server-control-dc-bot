# BUILD
FROM arm64v8/golang:latest AS build

WORKDIR /build

COPY . .

RUN go mod download

RUN make compile-arm64

#####################################

FROM arm64v8/alpine:3.14

WORKDIR /app

COPY --from=build /build/bin/main .

ENTRYPOINT ["/app/main"]