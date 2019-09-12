FROM golang:latest
LABEL maintainer="Payam Madjidi <pmadjidi@gmail.com>"
WORKDIR /app
COPY . .
RUN go get .
RUN go build --race -o main .
EXPOSE 8080
CMD ["./main"]