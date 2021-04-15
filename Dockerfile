FROM golang:1.16.3-buster

COPY build/bots /bots

ENTRYPOINT /bots
