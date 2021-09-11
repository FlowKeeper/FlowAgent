FROM alpine:latest
RUN mkdir /app
COPY /tmp/flowagent /app/agent

CMD ["/app/agent"]

