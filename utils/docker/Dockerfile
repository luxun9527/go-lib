FROM debian:stretch-slim

RUN mkdir "/app"
WORKDIR "/app"

COPY ./docker /app/docker
RUN chmod +x /app/docker
EXPOSE 8888
ENTRYPOINT ["./docker"]