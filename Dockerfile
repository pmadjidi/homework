FROM golang:latest
LABEL maintainer="Payam Madjidi <pmadjidi@gmail.com>"
ENV MAXQUEUELENGTH  200000
ENV	MAXITERATIONLIMIT 2000
ENV MAXNUMBEROFSTEPSINPUT 1000
ENV	MAXNUMBERSOFWALKERS  1000000
ENV	MAXNUMBEROFGROUPS  100000
ENV	MAXNUMBEROFWALKERSINGROUP 2000
ENV TIMEOUT 1
WORKDIR /app
COPY . .
RUN go get .
RUN go build --race -o main .
EXPOSE 8080
CMD ["./main"]