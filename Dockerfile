FROM golang:1.17-bullseye AS build
COPY . /src
WORKDIR /src
RUN ls -lah .
RUN CGO_ENABLED=0 go build -o /src/agent .

FROM alpine:3
RUN mkdir /app
COPY --from=build /src/agent /app/agent

CMD ["/app/agent"]

