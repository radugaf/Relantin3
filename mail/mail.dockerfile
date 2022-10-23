FROM alpine:latest

RUN mkdir /app

COPY mailService /app
COPY templates /templates


CMD ["/app/mailService"]