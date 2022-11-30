# Use AWS ECR public image
FROM public.ecr.aws/docker/library/golang:latest AS build

WORKDIR /Users/richa/GolandProjects/GoAgentProject

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o logglyassignment .

# Use AWS ECR public image
FROM public.ecr.aws/docker/library/alpine:latest

WORKDIR /GoAgentProject/

COPY --from=build /Users/richa/GolandProjects/GoAgentProject/logglyassignment .

ENV LOGGLY_TOKEN=02b996b5-0063-48c7-87e2-fbd03ffac6c0

ENV API_KEY=aqKf60jG6EIL9LqZKSf5KGnRAH1prGPe

CMD ["./logglyassignment"]
