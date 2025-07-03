FROM ubuntu:20.04
COPY littlebook /app/littlebook
WORKDIR /app
CMD ["/app/littlebook"]