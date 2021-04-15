FROM golang:1.16.0-buster

COPY build/bots /bots

ENTRYPOINT /bots
