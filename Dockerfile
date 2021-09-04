FROM golang:1.17-bullseye

COPY . /src
RUN cd /src && CGO_ENABLED=0 go build -o /src/agent .

FROM alpine:latest
RUN mkdir /app
COPY --from=0 /src/agent /app/agent

CMD ["/app/agent"]

