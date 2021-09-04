FROM alpine:3

RUN mkdir /app
RUN apk add --no-cache git make musl-dev go=1.16.7-r0
COPY . /src
RUN cd /src && go build -o /app/agent .
RUN apk del git make musl-dev go

CMD ["/app/agent"]
