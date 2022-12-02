FROM golang:latest AS build

WORKDIR /Users/richa/GolandProjects/GoAgentProject

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -o logglyassignment .

FROM public.ecr.aws/docker/library/alpine:latest

WORKDIR /GoAgentProject/

COPY --from=build /Users/richa/GolandProjects/GoAgentProject/logglyassignment ./

ENV LOGGLY_TOKEN=02b996b5-0063-48c7-87e2-fbd03ffac6c0

ENV API_KEY=aqKf60jG6EIL9LqZKSf5KGnRAH1prGPe

ENV AWS_DEFAULT_REGION=us-east-1

ENV AWS_ACCESS_KEY_ID=AKIA34XNLPJYL26ATFTT

ENV AWS_SECRET_ACCESS_KEY=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

CMD ["./logglyassignment", "eval $(AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION)","eval $(LOGGLY_TOKEN=$LOGGLY_TOKEN)","eval $(API_KEY=$API_KEY)","eval $(AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID)","eval $(AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY)"]