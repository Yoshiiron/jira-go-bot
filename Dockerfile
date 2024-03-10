FROM golang:1.22 AS Builder
WORKDIR /build
COPY ./ ./
RUN go build

FROM debian:bookworm-slim
COPY --from=Builder /build/jira-go-bot /usr/local/bin
RUN apt-get update && apt-get install -y ca-certificates && mkdir /DB
ENTRYPOINT ["jira-go-bot"]