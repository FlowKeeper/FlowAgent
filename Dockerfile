FROM golang:1.17-bullseye

COPY . /src
WORKDIR /src
RUN ls -lah .
RUN CGO_ENABLED=0 go build -o /src/agent .

FROM alpine:latest
RUN mkdir /app
COPY --from=0 /src/agent /app/agent

CMD ["/app/agent"]

