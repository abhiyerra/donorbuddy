FROM ubuntu:16.04

ADD . /donorbuddy

WORKDIR /

EXPOSE 8080

ENTRYPOINT ["/donorbuddy/donorbuddy", "/donorbuddy/config.dev.json"]
