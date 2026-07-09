FROM ubuntu:latest
LABEL authors="longer"

ENTRYPOINT ["top", "-b"]