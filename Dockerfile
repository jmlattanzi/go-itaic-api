FROM golang:latest
WORKDIR /usr/go
ADD . .
RUN
RUN dep ensure
EXPOSE 8000
CMD ["go", "run", "main.go"]