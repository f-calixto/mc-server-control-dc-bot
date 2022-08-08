FROM arm64v8/alpine:3.14

COPY ./bin .

CMD ["./main"]