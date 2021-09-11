FROM alpine:latest
RUN mkdir /app
COPY flowagent /app/agent

CMD ["/app/agent"]

