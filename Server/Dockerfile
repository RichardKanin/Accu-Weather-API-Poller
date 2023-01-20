FROM golang:latest AS build

WORKDIR /Users/richa/GolandProjects/http-project

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -o httpserver .

FROM alpine:latest

WORKDIR /http-project/

COPY --from=build /Users/richa/GolandProjects/http-project/httpserver ./

ENV LOGGLY_TOKEN=02b996b5-0063-48c7-87e2-fbd03ffac6c0

CMD ["./httpserver"]