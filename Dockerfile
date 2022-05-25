FROM alpine:3.16.0
WORKDIR /
COPY db db
COPY owl.linux owl

CMD ["/owl", "server", "--port", "8080"]

EXPOSE 8080
