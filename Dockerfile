# BUILD
FROM arm64v8/golang:latest AS build

WORKDIR /build

COPY . .

RUN make compile-arm64

#####################################

# RUN
FROM arm64v8/alpine:3.14

WORKDIR /app

COPY --from=build /build/bin/main .

ENTRYPOINT ["/app/main"]