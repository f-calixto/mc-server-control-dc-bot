FROM arm64v8/golang:latest AS build

COPY . .

RUN go mod install

RUN make compile-arm64

FROM arm64v8/alpine:3.14

WORKDIR /app

COPY --from=build ./bin .

CMD ["./main"]