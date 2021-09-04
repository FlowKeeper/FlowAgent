FROM golang:1.17-alpine

COPY . /src
RUN apk add gcc musl-dev
RUN cd /src && go build -o /src/agent .

FROM alpine:latest
RUN mkdir /app
COPY --from=0 /src/agent /app/agent

CMD ["/app/agent"]

